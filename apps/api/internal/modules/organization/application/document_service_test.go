package application_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepo struct {
	mock.Mock
}

func (m *MockDocumentRepo) CreateDocument(ctx context.Context, doc *domain.OrganizationDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentByID(ctx context.Context, id uuid.UUID) (*domain.OrganizationDocument, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrganizationDocument), args.Error(1)
}

func (m *MockDocumentRepo) ListDocumentsByOrganization(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationDocument, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.OrganizationDocument), args.Error(1)
}

func (m *MockDocumentRepo) GetMaxVersionForDocumentType(ctx context.Context, orgID uuid.UUID, docType string) (int, error) {
	args := m.Called(ctx, orgID, docType)
	return args.Int(0), args.Error(1)
}

func (m *MockDocumentRepo) UpdateDocumentReview(ctx context.Context, id uuid.UUID, status domain.OrganizationDocumentStatus, reviewerID uuid.UUID, rejectionReason *string) error {
	args := m.Called(ctx, id, status, reviewerID, rejectionReason)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetApprovedDocumentTypes(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	args := m.Called(ctx, key, r, contentType)
	return args.Error(0)
}

func (m *MockStorageService) PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *MockStorageService) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorageService) DeleteObject(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockStorageService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorageService) GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, expiry)
	return args.String(0), args.Error(1)
}

func TestDocumentService_GetDocumentForViewing_Unauthenticated(t *testing.T) {
	docRepo := new(MockDocumentRepo)
	orgRepo := new(MockOrgRepo)
	mockStorage := new(MockStorageService)

	svc := application.NewOrganizationDocumentApplicationService(
		&server.Server{},
		docRepo,
		orgRepo,
		mockStorage,
	)

	_, _, _, _, err := svc.GetDocumentForViewing(
		context.Background(),
		"", // empty callerID
		uuid.New(),
		uuid.New(),
		false,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication required")
}

func TestDocumentService_GetDocumentForViewing_MismatchOrg(t *testing.T) {
	docRepo := new(MockDocumentRepo)
	orgRepo := new(MockOrgRepo)
	mockStorage := new(MockStorageService)

	svc := application.NewOrganizationDocumentApplicationService(
		&server.Server{},
		docRepo,
		orgRepo,
		mockStorage,
	)

	docID := uuid.New()
	realOrgID := uuid.New()
	wrongOrgID := uuid.New()

	docRepo.On("GetDocumentByID", mock.Anything, docID).Return(&domain.OrganizationDocument{
		ID:               docID,
		OrganizationID:   realOrgID,
		DocumentType:     "registration_certificate",
		OriginalFilename: "cac.pdf",
		StorageKey:       "organizations/" + realOrgID.String() + "/cac.pdf",
		MIMEType:         "application/pdf",
		FileSizeBytes:    1024,
	}, nil)

	_, _, _, _, err := svc.GetDocumentForViewing(
		context.Background(),
		uuid.New().String(),
		wrongOrgID,
		docID,
		false,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Document does not belong to the requested organization")
}
