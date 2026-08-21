package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorageService implements ObjectStorage using local filesystem storage for dev/test environments.
type LocalStorageService struct {
	baseDir string
	baseURL string
}

func NewLocalStorageService(baseDir, baseURL string) (*LocalStorageService, error) {
	if baseDir == "" {
		baseDir = "./storage/documents"
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to initialize local storage directory: %w", err)
	}
	return &LocalStorageService{
		baseDir: baseDir,
		baseURL: baseURL,
	}, nil
}

func (s *LocalStorageService) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	return s.PutObject(ctx, key, r, 0, contentType)
}

func (s *LocalStorageService) PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	cleanKey := filepath.Clean(key)
	fullPath := filepath.Join(s.baseDir, cleanKey)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create storage object file: %w", err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("failed to write object bytes: %w", err)
	}
	if size > 0 && written != size {
		_ = os.Remove(fullPath)
		return fmt.Errorf("written byte count %d mismatch expected size %d", written, size)
	}

	return nil
}

func (s *LocalStorageService) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.GetObject(ctx, key)
}

func (s *LocalStorageService) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	cleanKey := filepath.Clean(key)
	fullPath := filepath.Join(s.baseDir, cleanKey)

	file, err := os.Open(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("object not found")
		}
		return nil, err
	}
	return file, nil
}

func (s *LocalStorageService) Delete(ctx context.Context, key string) error {
	return s.DeleteObject(ctx, key)
}

func (s *LocalStorageService) DeleteObject(ctx context.Context, key string) error {
	cleanKey := filepath.Clean(key)
	fullPath := filepath.Join(s.baseDir, cleanKey)
	err := os.Remove(fullPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (s *LocalStorageService) Exists(ctx context.Context, key string) (bool, error) {
	cleanKey := filepath.Clean(key)
	fullPath := filepath.Join(s.baseDir, cleanKey)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (s *LocalStorageService) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	cleanKey := filepath.ToSlash(filepath.Clean(key))
	if s.baseURL != "" {
		return fmt.Sprintf("%s/api/v1/storage/download?key=%s&expires=%d", strings.TrimSuffix(s.baseURL, "/"), cleanKey, time.Now().Add(expiration).Unix()), nil
	}
	return fmt.Sprintf("/api/v1/storage/download?key=%s&expires=%d", cleanKey, time.Now().Add(expiration).Unix()), nil
}
