package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

// S3StorageService provides high-performance, zero-daemon S3/Cloudflare R2 storage integration using AWS SigV4.
type S3StorageService struct {
	endpoint   string
	bucket     string
	region     string
	accessKey  string
	secretKey  string
	publicURL  string
	httpClient *http.Client
}

func NewS3StorageService(cfg StorageConfig) (*S3StorageService, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("S3_ENDPOINT is required when STORAGE_PROVIDER is s3/r2")
	}
	if cfg.Bucket == "" {
		return nil, errors.New("S3_BUCKET is required when STORAGE_PROVIDER is s3/r2")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, errors.New("S3_ACCESS_KEY and S3_SECRET_KEY are required when STORAGE_PROVIDER is s3/r2")
	}

	region := cfg.Region
	if region == "" {
		region = "auto" // Cloudflare R2 default region
	}

	endpoint := strings.TrimRight(cfg.Endpoint, "/")

	return &S3StorageService{
		endpoint:   endpoint,
		bucket:     cfg.Bucket,
		region:     region,
		accessKey:  cfg.AccessKey,
		secretKey:  cfg.SecretKey,
		publicURL:  strings.TrimRight(cfg.PublicURL, "/"),
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (s *S3StorageService) buildObjectURI(key string) string {
	cleanKey := strings.TrimPrefix(path.Clean("/"+key), "/")
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, cleanKey)
}

func (s *S3StorageService) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	return s.PutObject(ctx, key, r, 0, contentType)
}

func (s *S3StorageService) PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	rawURL := s.buildObjectURI(key)
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid storage object URL: %w", err)
	}

	var body []byte
	if reader != nil {
		body, err = io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read upload body: %w", err)
		}
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	s.signRequest(req, body)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send object to S3/R2 storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3/R2 put failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *S3StorageService) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.GetObject(ctx, key)
}

func (s *S3StorageService) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	rawURL := s.buildObjectURI(key)
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid storage object URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get request: %w", err)
	}

	s.signRequest(req, nil)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object from S3/R2: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, errors.New("object not found")
	}

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("S3/R2 get failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return resp.Body, nil
}

func (s *S3StorageService) Delete(ctx context.Context, key string) error {
	return s.DeleteObject(ctx, key)
}

func (s *S3StorageService) DeleteObject(ctx context.Context, key string) error {
	rawURL := s.buildObjectURI(key)
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid storage object URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	s.signRequest(req, nil)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete object in S3/R2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3/R2 delete failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (s *S3StorageService) Exists(ctx context.Context, key string) (bool, error) {
	rawURL := s.buildObjectURI(key)
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("invalid storage object URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("failed to create head request: %w", err)
	}

	s.signRequest(req, nil)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("head request returned unexpected status %d", resp.StatusCode)
}

func (s *S3StorageService) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	if expiration <= 0 {
		expiration = 15 * time.Minute
	}

	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, s.region)
	cleanKey := strings.TrimPrefix(path.Clean("/"+key), "/")
	objectPath := fmt.Sprintf("/%s/%s", s.bucket, cleanKey)

	u, err := url.Parse(s.endpoint)
	if err != nil {
		return "", err
	}

	canonicalURI := objectPath
	expiresSeconds := int(expiration.Seconds())

	queryParams := url.Values{}
	queryParams.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	queryParams.Set("X-Amz-Credential", fmt.Sprintf("%s/%s", s.accessKey, credentialScope))
	queryParams.Set("X-Amz-Date", amzDate)
	queryParams.Set("X-Amz-Expires", fmt.Sprintf("%d", expiresSeconds))
	queryParams.Set("X-Amz-SignedHeaders", "host")

	canonicalQueryString := queryParams.Encode()

	canonicalHeaders := fmt.Sprintf("host:%s\n", u.Host)
	signedHeaders := "host"
	payloadHash := "UNSIGNED-PAYLOAD"

	canonicalRequest := strings.Join([]string{
		http.MethodGet,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	hReq := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hex.EncodeToString(hReq[:]),
	}, "\n")

	signingKey := getSignatureKey(s.secretKey, dateStamp, s.region, "s3")
	signature := hmacSHA256(signingKey, []byte(stringToSign))
	queryParams.Set("X-Amz-Signature", hex.EncodeToString(signature))

	return fmt.Sprintf("%s%s?%s", s.endpoint, objectPath, queryParams.Encode()), nil
}

// AWS SigV4 Signing Helpers

func (s *S3StorageService) signRequest(req *http.Request, body []byte) {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")

	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("Host", req.URL.Host)

	var payloadHash string
	if len(body) == 0 {
		emptyHash := sha256.Sum256([]byte(""))
		payloadHash = hex.EncodeToString(emptyHash[:])
	} else {
		hash := sha256.Sum256(body)
		payloadHash = hex.EncodeToString(hash[:])
	}
	req.Header.Set("x-amz-content-sha256", payloadHash)

	// Canonical Headers
	headers := []string{"host", "x-amz-content-sha256", "x-amz-date"}
	if ct := req.Header.Get("Content-Type"); ct != "" {
		headers = append(headers, "content-type")
	}
	sort.Strings(headers)

	var canonHeaders strings.Builder
	for _, h := range headers {
		canonHeaders.WriteString(fmt.Sprintf("%s:%s\n", h, strings.TrimSpace(req.Header.Get(h))))
	}
	signedHeaders := strings.Join(headers, ";")

	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		canonHeaders.String(),
		signedHeaders,
		payloadHash,
	}, "\n")

	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, s.region)
	hReq := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hex.EncodeToString(hReq[:]),
	}, "\n")

	signingKey := getSignatureKey(s.secretKey, dateStamp, s.region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	authHeader := fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey,
		credentialScope,
		signedHeaders,
		signature,
	)
	req.Header.Set("Authorization", authHeader)
}

func hmacSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(regionName))
	kService := hmacSHA256(kRegion, []byte(serviceName))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}
