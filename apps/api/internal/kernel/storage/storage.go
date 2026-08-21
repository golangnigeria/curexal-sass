package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// ObjectStorage is the canonical vendor-independent abstraction for object storage.
type ObjectStorage interface {
	Put(ctx context.Context, key string, r io.Reader, contentType string) error
	PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	DeleteObject(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}

// ObjectStorageService is an alias for backward compatibility across all modules.
type ObjectStorageService = ObjectStorage

// StorageConfig holds provider configuration settings.
type StorageConfig struct {
	Provider  string // "local", "s3", "r2"
	BaseDir   string // Local directory (default: "./storage/documents")
	Endpoint  string // S3 / R2 Endpoint (e.g. https://<accountid>.r2.cloudflarestorage.com)
	Bucket    string // S3 / R2 Bucket name
	Region    string // S3 / R2 Region (default: "auto")
	AccessKey string // S3 / R2 Access Key ID
	SecretKey string // S3 / R2 Secret Access Key
	PublicURL string // Optional custom domain / CDN URL
}

// New initializes the appropriate ObjectStorage implementation based on configuration.
func New(cfg StorageConfig, baseURL string) (ObjectStorage, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	switch provider {
	case "s3", "r2", "cloudflare_r2":
		return NewS3StorageService(cfg)
	case "local", "":
		return NewLocalStorageService(cfg.BaseDir, baseURL)
	default:
		return nil, fmt.Errorf("unsupported storage provider %q (supported: local, s3, r2)", cfg.Provider)
	}
}

// Security Validation Utilities

var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var AllowedMIMETypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
}

func ValidateExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !AllowedExtensions[ext] {
		return fmt.Errorf("unsupported file extension %q. Allowed: .pdf, .jpg, .jpeg, .png", ext)
	}
	return nil
}

func ValidateMIMEType(content []byte) (string, error) {
	if len(content) == 0 {
		return "", errors.New("empty file content")
	}
	mime := http.DetectContentType(content)
	if idx := strings.Index(mime, ";"); idx != -1 {
		mime = strings.TrimSpace(mime[:idx])
	}
	if !AllowedMIMETypes[mime] {
		return "", fmt.Errorf("unsupported MIME type %q. Allowed: application/pdf, image/jpeg, image/png", mime)
	}
	return mime, nil
}

func CalculateSHA256(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

func BuildDeterministicStorageKey(orgID, docType, docID string, version int) string {
	cleanOrg := filepath.Clean(orgID)
	cleanDocType := filepath.Clean(docType)
	cleanDocID := filepath.Clean(docID)
	return fmt.Sprintf("organizations/%s/documents/%s/%s/v%d", cleanOrg, cleanDocType, cleanDocID, version)
}

func NewMemoryBufferReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}
