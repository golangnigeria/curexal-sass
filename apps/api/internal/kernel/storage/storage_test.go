package storage_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/golangnigeria/curexal/internal/kernel/storage"
)

func TestStorageSecurityValidations(t *testing.T) {
	t.Run("Validate Extension Allowed", func(t *testing.T) {
		err := storage.ValidateExtension("operating_license.pdf")
		if err != nil {
			t.Fatalf("expected nil error for valid pdf, got %v", err)
		}
	})

	t.Run("Validate Extension Rejected", func(t *testing.T) {
		err := storage.ValidateExtension("malicious_script.exe")
		if err == nil {
			t.Fatalf("expected error for exe extension, got nil")
		}
	})

	t.Run("Validate MIME Type Detection", func(t *testing.T) {
		pdfHeader := []byte("%PDF-1.5 test content")
		mime, err := storage.ValidateMIMEType(pdfHeader)
		if err != nil {
			t.Fatalf("expected valid mime type, got error: %v", err)
		}
		if mime != "application/pdf" {
			t.Fatalf("expected application/pdf, got %s", mime)
		}
	})

	t.Run("Calculate SHA-256 Checksum", func(t *testing.T) {
		content := []byte("curexal security document")
		checksum := storage.CalculateSHA256(content)
		if checksum == "" || len(checksum) != 64 {
			t.Fatalf("expected 64-char sha256 hex string, got %s", checksum)
		}
	})

	t.Run("Deterministic Storage Key", func(t *testing.T) {
		key := storage.BuildDeterministicStorageKey("org-123", "operating_license", "doc-999", 2)
		if !strings.Contains(key, "organizations/org-123/documents/operating_license/doc-999/v2") {
			t.Fatalf("unexpected storage key format: %s", key)
		}
	})
}

func TestLocalStorageService(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "curexal_storage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var s storage.ObjectStorage
	s, err = storage.New(storage.StorageConfig{
		Provider: "local",
		BaseDir:  tempDir,
	}, "http://localhost:8080")
	if err != nil {
		t.Fatalf("failed to create storage service: %v", err)
	}

	ctx := context.Background()
	key := "test/file.pdf"
	content := []byte("sample document payload")

	// Exists Before Put
	exists, errExists := s.Exists(ctx, key)
	if errExists != nil || exists {
		t.Fatalf("expected key not to exist initially")
	}

	// Put Object
	errPut := s.Put(ctx, key, storage.NewMemoryBufferReader(content), "application/pdf")
	if errPut != nil {
		t.Fatalf("Put failed: %v", errPut)
	}

	// Exists After Put
	exists, errExists = s.Exists(ctx, key)
	if errExists != nil || !exists {
		t.Fatalf("expected key to exist after put")
	}

	// Presigned URL
	url, errURL := s.GeneratePresignedURL(ctx, key, 0)
	if errURL != nil || !strings.Contains(url, "key=test/file.pdf") {
		t.Fatalf("GeneratePresignedURL failed or unexpected: %s", url)
	}

	// Get Object
	rc, errGet := s.Get(ctx, key)
	if errGet != nil {
		t.Fatalf("Get failed: %v", errGet)
	}
	rc.Close()

	// Delete Object
	errDel := s.Delete(ctx, key)
	if errDel != nil {
		t.Fatalf("Delete failed: %v", errDel)
	}

	// Exists After Delete
	exists, _ = s.Exists(ctx, key)
	if exists {
		t.Fatalf("expected key not to exist after delete")
	}
}

func TestS3ConfigValidation(t *testing.T) {
	t.Run("Missing Endpoint", func(t *testing.T) {
		_, err := storage.New(storage.StorageConfig{
			Provider:  "s3",
			Endpoint:  "",
			Bucket:    "curexal-docs",
			AccessKey: "key",
			SecretKey: "secret",
		}, "")
		if err == nil {
			t.Fatalf("expected error when endpoint missing")
		}
	})

	t.Run("Valid S3 Service Creation", func(t *testing.T) {
		svc, err := storage.New(storage.StorageConfig{
			Provider:  "s3",
			Endpoint:  "https://example.r2.cloudflarestorage.com",
			Bucket:    "curexal-docs",
			AccessKey: "key",
			SecretKey: "secret",
			Region:    "auto",
		}, "")
		if err != nil {
			t.Fatalf("expected valid s3 storage creation, got: %v", err)
		}
		if svc == nil {
			t.Fatalf("expected non-nil service")
		}
	})
}
