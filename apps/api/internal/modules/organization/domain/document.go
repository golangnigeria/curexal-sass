package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrganizationDocumentStatus string

const (
	DocumentStatusPending  OrganizationDocumentStatus = "pending"
	DocumentStatusApproved OrganizationDocumentStatus = "approved"
	DocumentStatusRejected OrganizationDocumentStatus = "rejected"
	DocumentStatusExpired  OrganizationDocumentStatus = "expired"
)

type OrganizationDocument struct {
	ID               uuid.UUID                  `json:"id"`
	OrganizationID   uuid.UUID                  `json:"organizationId"`
	DocumentType     string                     `json:"documentType"`
	OriginalFilename string                     `json:"originalFilename"`
	StorageKey       string                     `json:"storageKey"`
	MIMEType         string                     `json:"mimeType"`
	FileSizeBytes    int64                      `json:"fileSizeBytes"`
	ChecksumSHA256   string                     `json:"checksumSha256"`
	UploadedBy       uuid.UUID                  `json:"uploadedBy"`
	UploadedAt       time.Time                  `json:"uploadedAt"`
	Status           OrganizationDocumentStatus `json:"status"`
	Version          int                        `json:"version"`
	ReviewedBy       *uuid.UUID                 `json:"reviewedBy,omitempty"`
	ReviewedAt       *time.Time                 `json:"reviewedAt,omitempty"`
	RejectionReason  *string                    `json:"rejectionReason,omitempty"`
	CreatedAt        time.Time                  `json:"createdAt"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, doc *OrganizationDocument) error
	GetDocumentByID(ctx context.Context, id uuid.UUID) (*OrganizationDocument, error)
	ListDocumentsByOrganization(ctx context.Context, orgID uuid.UUID) ([]OrganizationDocument, error)
	GetMaxVersionForDocumentType(ctx context.Context, orgID uuid.UUID, docType string) (int, error)
	UpdateDocumentReview(ctx context.Context, id uuid.UUID, status OrganizationDocumentStatus, reviewerID uuid.UUID, rejectionReason *string) error
	GetApprovedDocumentTypes(ctx context.Context, orgID uuid.UUID) ([]string, error)
}
