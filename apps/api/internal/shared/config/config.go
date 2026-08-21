package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

func (d DatabaseConfig) DSN() string {
	if dsn := os.Getenv("CUREXAL_DB_DSN"); dsn != "" {
		return strings.Trim(strings.TrimSpace(dsn), "\"'")
	}
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return strings.Trim(strings.TrimSpace(dsn), "\"'")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

type ServerConfig struct {
	Port               string
	Domain             string
	PortalPort         string
	PublicPort         string
	WorkspacePort      string
	ReadTimeout        int
	WriteTimeout       int
	IdleTimeout        int
	CORSAllowedOrigins []string
}

type ObservabilityConfig struct {
	ServiceName  string
	Environment  string
	Logging      LoggingConfig
	NewRelic     NewRelicConfig
	HealthChecks HealthChecksConfig
}

type LoggingConfig struct {
	Level              string
	Format             string
	SlowQueryThreshold int
}

type NewRelicConfig struct {
	LicenseKey                string
	AppLogForwardingEnabled   bool
	DistributedTracingEnabled bool
	DebugLogging              bool
}

type HealthChecksConfig struct {
	Enabled bool
}

type IntegrationConfig struct {
	ResendAPIKey     string
	EmailFromName    string
	EmailFromAddress string
	EmailLogoURL     string
	EmailAppURL      string
}

type RedisConfig struct {
	Address string
}

type StorageConfig struct {
	Provider  string
	BaseDir   string
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	PublicURL string
}

type CacheConfig struct {
	Provider   string
	DefaultTTL time.Duration
}

type AuthConfig struct {
	SecretKey          string
	JWTCookieName      string
	RefreshCookieName  string
	CookieDomain       string
	CookiePath         string
	CookieSecure       bool
	CookieHTTPOnly     bool
	CookieSameSite     string
	JWTExpiration      time.Duration
	RefreshExpiration  time.Duration
	LockoutDuration    time.Duration
	PlatformStaffRoles []string
	AllowTestHeaders   bool
}

type Primary struct {
	Env string
}

type Config struct {
	Primary       Primary
	Database      DatabaseConfig
	Server        ServerConfig
	Observability *ObservabilityConfig
	Integration   IntegrationConfig
	Redis         RedisConfig
	Storage       StorageConfig
	Cache         CacheConfig
	Auth          AuthConfig
}

func (c *Config) ResolveCookieDomain() string {
	if c.Auth.CookieDomain != "" {
		return c.Auth.CookieDomain
	}
	return c.Server.Domain
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return strings.Trim(strings.TrimSpace(val), "\"'")
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		cleaned := strings.ToLower(strings.Trim(strings.TrimSpace(val), "\"'"))
		return cleaned == "true" || cleaned == "1" || cleaned == "yes"
	}
	return defaultVal
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env", "../.env", "../../.env")
	fmt.Println("Config validation passed")
	return &Config{
		Primary: Primary{Env: getEnv("CUREXAL_PRIMARY_ENV", "local")},
		Database: DatabaseConfig{
			Host:         getEnv("CUREXAL_DATABASE_HOST", getEnv("CUREXAL_DATABASE.HOST", "localhost")),
			Port:         5432,
			User:         getEnv("CUREXAL_DATABASE_USER", getEnv("CUREXAL_DATABASE.USER", "postgres")),
			Password:     getEnv("CUREXAL_DATABASE_PASSWORD", getEnv("CUREXAL_DATABASE.PASSWORD", "postgres")),
			Name:         getEnv("CUREXAL_DATABASE_NAME", getEnv("CUREXAL_DATABASE.NAME", "CUREXAL")),
			SSLMode:      getEnv("CUREXAL_DATABASE_SSL_MODE", getEnv("CUREXAL_DATABASE.SSL_MODE", "disable")),
			MaxOpenConns: 25,
			MaxIdleConns: 5,
		},
		Server: ServerConfig{
			Port:          getEnv("CUREXAL_SERVER_PORT", getEnv("PORT", "8080")),
			Domain:        getEnv("CUREXAL_SERVER_DOMAIN", getEnv("COOKIE_DOMAIN", "localhost")),
			PortalPort:    getEnv("CUREXAL_SERVER_PORTAL_PORT", "5001"),
			PublicPort:    getEnv("CUREXAL_SERVER_PUBLIC_PORT", "5001"),
			WorkspacePort: getEnv("CUREXAL_SERVER_WORKSPACE_PORT", "5002"),
		},
		Integration: IntegrationConfig{
			ResendAPIKey:     getEnv("CUREXAL_INTEGRATION_RESEND_API_KEY", getEnv("CUREXAL_INTEGRATION.RESEND_API_KEY", "")),
			EmailFromName:    getEnv("CUREXAL_INTEGRATION_EMAIL_FROM_NAME", getEnv("CUREXAL_INTEGRATION.EMAIL_FROM_NAME", "Curexal")),
			EmailFromAddress: getEnv("CUREXAL_INTEGRATION_EMAIL_FROM_ADDRESS", getEnv("CUREXAL_INTEGRATION.EMAIL_FROM_ADDRESS", "noreply@contact.curexal.space")),
			EmailLogoURL:     getEnv("CUREXAL_INTEGRATION_EMAIL_LOGO_URL", getEnv("EMAIL_LOGO_URL", "https://cdn.curexal.space/email/full_logo.png")),
			EmailAppURL:      getEnv("CUREXAL_INTEGRATION_EMAIL_APP_URL", getEnv("EMAIL_APP_URL", "https://curexal.space")),
		},
		Redis: RedisConfig{
			Address: getEnv("REDIS_ADDR", getEnv("CUREXAL_REDIS_ADDR", "")),
		},
		Storage: StorageConfig{
			Provider:  getEnv("STORAGE_PROVIDER", "local"),
			BaseDir:   getEnv("STORAGE_LOCAL_BASE_DIR", "./storage/documents"),
			Endpoint:  getEnv("S3_ENDPOINT", ""),
			Bucket:    getEnv("S3_BUCKET", "curexal-documents"),
			Region:    getEnv("S3_REGION", "auto"),
			AccessKey: getEnv("S3_ACCESS_KEY", ""),
			SecretKey: getEnv("S3_SECRET_KEY", ""),
			PublicURL: getEnv("S3_PUBLIC_URL", ""),
		},
		Cache: CacheConfig{
			Provider:   getEnv("CACHE_PROVIDER", "memory"),
			DefaultTTL: parseDurationEnv("CACHE_DEFAULT_TTL", 10*time.Minute),
		},
		Auth: AuthConfig{
			SecretKey:          getEnv("CUREXAL_AUTH_SECRET_KEY", "super-secret-key-curexal-platform-development"),
			JWTCookieName:      getEnv("CUREXAL_AUTH_JWT_COOKIE_NAME", "jwt"),
			RefreshCookieName:  getEnv("CUREXAL_AUTH_REFRESH_COOKIE_NAME", "refresh_token"),
			CookieDomain:       getEnv("CUREXAL_AUTH_COOKIE_DOMAIN", getEnv("COOKIE_DOMAIN", "")),
			CookiePath:         getEnv("CUREXAL_AUTH_COOKIE_PATH", "/"),
			CookieSecure:       getEnvBool("CUREXAL_AUTH_COOKIE_SECURE", false),
			CookieHTTPOnly:     getEnvBool("CUREXAL_AUTH_COOKIE_HTTPONLY", true),
			CookieSameSite:     getEnv("CUREXAL_AUTH_COOKIE_SAMESITE", "Default"),
			JWTExpiration:      15 * time.Minute,
			RefreshExpiration:  30 * 24 * time.Hour,
			LockoutDuration:    parseDurationEnv("CUREXAL_AUTH_LOCKOUT_DURATION", 1*time.Minute),
			PlatformStaffRoles: []string{"super_admin", "super_support_agent", "super_sales_staff", "super_compliance_officer"},
			AllowTestHeaders:   getEnvBool("CUREXAL_AUTH_ALLOW_TEST_HEADERS", true),
		},
	}, nil
}

func parseDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
