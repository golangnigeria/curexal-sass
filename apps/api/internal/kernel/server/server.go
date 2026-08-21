package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/cache"
	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/kernel/storage"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/logger"
	"github.com/golangnigeria/curexal/internal/shared/mailer"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type JobQueueClient struct{}

func (c *JobQueueClient) Enqueue(_ interface{}) (interface{}, error) {
	return nil, nil
}

func (c *JobQueueClient) EnqueueContext(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

type JobClient struct {
	Client *JobQueueClient
}

func (j *JobClient) EnqueueEmail(_ interface{}) error { return nil }

// Server is the platform HTTP server wrapper and dependency holder.
type Server struct {
	Config     *config.Config
	DB         *database.Database
	Logger     *zerolog.Logger
	Redis      *redis.Client
	Cache      cache.CacheService
	Storage    storage.ObjectStorage
	Job        *JobClient
	Mailer     *mailer.Mailer
	Echo       *echo.Echo
	httpServer *http.Server
}

// New initializes a new platform Server instance.
func New(cfg *config.Config, log *zerolog.Logger, _ *logger.LoggerService) (*Server, error) {
	db, err := database.New(cfg, log, nil)
	if err != nil {
		return nil, err
	}

	baseURL := ""
	if cfg.Server.Port != "" {
		baseURL = fmt.Sprintf("http://localhost:%s", cfg.Server.Port)
	}
	storageSvc, errStorage := storage.New(storage.StorageConfig{
		Provider:  cfg.Storage.Provider,
		BaseDir:   cfg.Storage.BaseDir,
		Endpoint:  cfg.Storage.Endpoint,
		Bucket:    cfg.Storage.Bucket,
		Region:    cfg.Storage.Region,
		AccessKey: cfg.Storage.AccessKey,
		SecretKey: cfg.Storage.SecretKey,
		PublicURL: cfg.Storage.PublicURL,
	}, baseURL)
	if errStorage != nil {
		log.Warn().Err(errStorage).Msg("storage service initialized with warning, falling back to local driver")
		storageSvc, _ = storage.NewLocalStorageService(cfg.Storage.BaseDir, baseURL)
	}

	cacheSvc := cache.New(cache.Config{
		Provider:     cfg.Cache.Provider,
		RedisAddress: cfg.Redis.Address,
		DefaultTTL:   cfg.Cache.DefaultTTL,
	})

	m := mailer.New(
		cfg.Integration.ResendAPIKey,
		cfg.Integration.EmailFromName,
		cfg.Integration.EmailFromAddress,
		log,
	)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusBadRequest
		msg := err.Error()
		errCode := "BAD_REQUEST"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if he.Message != nil {
				msg = fmt.Sprintf("%v", he.Message)
			}
			errCode = http.StatusText(code)
		} else if appErr, ok := err.(*errs.AppError); ok {
			code = appErr.HTTPStatus
			msg = appErr.Message
			errCode = appErr.Code
		} else {
			loweredMsg := strings.ToLower(msg)
			if strings.Contains(loweredMsg, "not found") {
				code = http.StatusNotFound
				errCode = "NOT_FOUND"
			} else if strings.Contains(loweredMsg, "unauthorized") || strings.Contains(loweredMsg, "unauthenticated") {
				code = http.StatusUnauthorized
				errCode = "UNAUTHORIZED"
			} else if strings.Contains(loweredMsg, "forbidden") || strings.Contains(loweredMsg, "permission") {
				code = http.StatusForbidden
				errCode = "FORBIDDEN"
			} else if strings.Contains(loweredMsg, "already taken") || strings.Contains(loweredMsg, "conflict") || strings.Contains(loweredMsg, "exists") {
				code = http.StatusConflict
				errCode = "CONFLICT"
			}
		}

		if code == http.StatusInternalServerError && msg == "" {
			msg = "An unexpected error occurred"
			errCode = "INTERNAL_SERVER_ERROR"
		}

		_ = c.JSON(code, map[string]interface{}{
			"message": msg,
			"data":    map[string]interface{}{},
			"meta": map[string]interface{}{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			},
			"links": map[string]interface{}{},
			"errors": []map[string]interface{}{
				{
					"code":    errCode,
					"message": msg,
				},
			},
		})
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			if origin == "" {
				return true, nil
			}
			// Allow any localhost / 127.0.0.1 origin or curexal domain across ports & subdomains
			if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") || strings.Contains(origin, "curexal.space") || strings.Contains(origin, "curexal.com") || strings.Contains(origin, "curexal.internal") {
				return true, nil
			}
			return true, nil
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Tenant-Slug", "X-Active-Tenant-ID", "X-Tenant-ID", "X-Access-Token", "X-User-ID", "X-User-Role", "X-Requested-With", "X-CSRF-Token"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Zerolog HTTP Request Logging Middleware for terminal debugging
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			l := log.Info()
			if err != nil {
				l = log.Error().Err(err)
			} else if res.Status >= 500 {
				l = log.Error()
			} else if res.Status >= 400 {
				l = log.Warn()
			}

			l.Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Str("ip", c.RealIP()).
				Str("latency", stop.Sub(start).String()).
				Msg("HTTP")

			return err
		}
	})

	return &Server{
		Config:  cfg,
		DB:      db,
		Logger:  log,
		Storage: storageSvc,
		Cache:   cacheSvc,
		Job:     &JobClient{Client: &JobQueueClient{}},
		Mailer:  m,
		Echo:    e,
	}, nil
}

func (s *Server) Start() error {
	port := s.Config.Server.Port
	if port == "" {
		port = "8080"
	}
	addr := port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	s.Logger.Info().
		Str("port", port).
		Str("url", "http://localhost"+addr).
		Msg("HTTP server listening")

	if s.Echo != nil {
		s.httpServer = &http.Server{
			Addr:    addr,
			Handler: s.Echo,
		}
		return s.httpServer.ListenAndServe()
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		_ = s.httpServer.Shutdown(ctx)
	}
	if s.DB != nil && s.DB.Pool != nil {
		s.DB.Pool.Close()
	}
	return nil
}

func (s *Server) SetupHTTPServer(_ interface{}) {}

// ExecuteInTenantTx executes a database function within a transaction scoped strictly to tenantSchema (SET LOCAL search_path = tenantSchema, public).
func (s *Server) ExecuteInTenantTx(ctx context.Context, tenantSchema string, fn func(tx pgx.Tx) error) error {
	if s.DB == nil || s.DB.Pool == nil {
		return nil
	}

	tx, err := s.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Connection search_path set strictly for this transaction ONLY
	setPathSQL := fmt.Sprintf("SET LOCAL search_path = %s, public", pgx.Identifier{tenantSchema}.Sanitize())
	if _, err := tx.Exec(ctx, setPathSQL); err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
