package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DocumentRepository struct {
	server *server.Server
}

func NewDocumentRepository(server *server.Server) *DocumentRepository {
	return &DocumentRepository{server: server}
}

func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *domain.OrganizationDocument) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		INSERT INTO organization.organization_documents (
			id, organization_id, document_type, original_filename, storage_key,
			mime_type, file_size_bytes, checksum_sha256, uploaded_by, uploaded_at,
			status, version
		) VALUES (
			@id, @org_id, @doc_type, @orig_filename, @storage_key,
			@mime_type, @size, @checksum, @uploaded_by, CURRENT_TIMESTAMP,
			@status, @version
		)
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{
		"id":            doc.ID.String(),
		"org_id":        doc.OrganizationID.String(),
		"doc_type":      doc.DocumentType,
		"orig_filename": doc.OriginalFilename,
		"storage_key":   doc.StorageKey,
		"mime_type":     doc.MIMEType,
		"size":          doc.FileSizeBytes,
		"checksum":      doc.ChecksumSHA256,
		"uploaded_by":   doc.UploadedBy.String(),
		"status":        string(doc.Status),
		"version":       doc.Version,
	})
	if err != nil {
		return fmt.Errorf("failed to insert organization document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) GetDocumentByID(ctx context.Context, id uuid.UUID) (*domain.OrganizationDocument, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT 
			id, organization_id, document_type, original_filename, storage_key,
			mime_type, file_size_bytes, checksum_sha256, uploaded_by, uploaded_at,
			status, version, reviewed_by, reviewed_at, rejection_reason,
			created_at, updated_at
		FROM organization.organization_documents
		WHERE id = @id
	`
	row := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"id": id.String()})
	d := &domain.OrganizationDocument{}
	var statusStr string
	var reviewedByStr *string

	err := row.Scan(
		&d.ID, &d.OrganizationID, &d.DocumentType, &d.OriginalFilename, &d.StorageKey,
		&d.MIMEType, &d.FileSizeBytes, &d.ChecksumSHA256, &d.UploadedBy, &d.UploadedAt,
		&statusStr, &d.Version, &reviewedByStr, &d.ReviewedAt, &d.RejectionReason,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("document not found")
		}
		return nil, fmt.Errorf("failed to query document: %w", err)
	}
	d.Status = domain.OrganizationDocumentStatus(statusStr)
	if reviewedByStr != nil {
		if uID, errParse := uuid.Parse(*reviewedByStr); errParse == nil {
			d.ReviewedBy = &uID
		}
	}
	return d, nil
}

func (r *DocumentRepository) ListDocumentsByOrganization(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationDocument, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT 
			id, organization_id, document_type, original_filename, storage_key,
			mime_type, file_size_bytes, checksum_sha256, uploaded_by, uploaded_at,
			status, version, reviewed_by, reviewed_at, rejection_reason,
			created_at, updated_at
		FROM organization.organization_documents
		WHERE organization_id = @org_id
		ORDER BY document_type ASC, version DESC
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"org_id": orgID.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to query organization documents: %w", err)
	}
	defer rows.Close()

	var docs []domain.OrganizationDocument
	for rows.Next() {
		d := domain.OrganizationDocument{}
		var statusStr string
		var reviewedByStr *string
		if err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.DocumentType, &d.OriginalFilename, &d.StorageKey,
			&d.MIMEType, &d.FileSizeBytes, &d.ChecksumSHA256, &d.UploadedBy, &d.UploadedAt,
			&statusStr, &d.Version, &reviewedByStr, &d.ReviewedAt, &d.RejectionReason,
			&d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan document row: %w", err)
		}
		d.Status = domain.OrganizationDocumentStatus(statusStr)
		if reviewedByStr != nil {
			if uID, errParse := uuid.Parse(*reviewedByStr); errParse == nil {
				d.ReviewedBy = &uID
			}
		}
		docs = append(docs, d)
	}
	if docs == nil {
		docs = []domain.OrganizationDocument{}
	}
	return docs, nil
}

func (r *DocumentRepository) GetMaxVersionForDocumentType(ctx context.Context, orgID uuid.UUID, docType string) (int, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT COALESCE(MAX(version), 0) 
		FROM organization.organization_documents 
		WHERE organization_id = @org_id AND document_type = @doc_type
	`
	var maxVer int
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"org_id": orgID.String(), "doc_type": docType}).Scan(&maxVer)
	if err != nil {
		return 0, fmt.Errorf("failed to query max document version: %w", err)
	}
	return maxVer, nil
}

func (r *DocumentRepository) UpdateDocumentReview(ctx context.Context, id uuid.UUID, status domain.OrganizationDocumentStatus, reviewerID uuid.UUID, rejectionReason *string) error {
	dbExec := r.server.DB.Conn(ctx)
	now := time.Now()
	stmt := `
		UPDATE organization.organization_documents
		SET status = @status, reviewed_by = @reviewer_id, reviewed_at = @now, rejection_reason = @reason, updated_at = CURRENT_TIMESTAMP
		WHERE id = @id
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{
		"status":      string(status),
		"reviewer_id": reviewerID.String(),
		"now":         now,
		"reason":      rejectionReason,
		"id":          id.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to update document review status: %w", err)
	}
	return nil
}

func (r *DocumentRepository) GetApprovedDocumentTypes(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT DISTINCT document_type 
		FROM organization.organization_documents 
		WHERE organization_id = @org_id AND status = 'approved'
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"org_id": orgID.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to query approved document types: %w", err)
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var dt string
		if err := rows.Scan(&dt); err == nil {
			types = append(types, dt)
		}
	}
	return types, nil
}
