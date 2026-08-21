package application

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/kernel/storage"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/google/uuid"
)

type DocumentWithPresignedURL struct {
	Document     domain.OrganizationDocument `json:"document"`
	PresignedURL string                      `json:"presignedUrl"`
}

type OrganizationDocumentApplicationService struct {
	server         *server.Server
	docRepo        domain.DocumentRepository
	orgRepo        domain.OrganizationRepository
	storageService storage.ObjectStorageService
}

func NewOrganizationDocumentApplicationService(
	server *server.Server,
	docRepo domain.DocumentRepository,
	orgRepo domain.OrganizationRepository,
	storageService storage.ObjectStorageService,
) *OrganizationDocumentApplicationService {
	return &OrganizationDocumentApplicationService{
		server:         server,
		docRepo:        docRepo,
		orgRepo:        orgRepo,
		storageService: storageService,
	}
}

// Category to Required Document Catalog Rules
var RequiredDocumentsPerCategory = map[string][]string{
	"laboratory":  {"operating_license", "registration_certificate"},
	"clinic":      {"operating_license", "facility_license"},
	"pharmacy":    {"operating_license", "pharmacy_license"},
	"hospital":    {"operating_license", "facility_license", "registration_certificate"},
	"radiology":   {"operating_license", "accreditation_certificate"},
	"healthcare":  {"operating_license"},
}

func (s *OrganizationDocumentApplicationService) UploadDocument(
	ctx context.Context,
	callerID string,
	orgID uuid.UUID,
	docType string,
	filename string,
	content []byte,
) (*domain.OrganizationDocument, error) {
	if callerID == "" {
		return nil, errs.NewUnauthorizedError("authentication required")
	}

	// 1. Verify caller membership in target organization
	var isMember bool
	errMember := s.server.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships 
			WHERE user_id = $1::uuid AND organization_id = $2::uuid AND is_active = TRUE
		)
	`, callerID, orgID.String()).Scan(&isMember)
	if errMember != nil || !isMember {
		// Also check platform admin bypass
		var isPAdmin bool
		_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, callerID).Scan(&isPAdmin)
		if !isPAdmin {
			return nil, errs.NewForbiddenError("User is not a member of the target organization")
		}
	}

	// 2. Validate File Extension
	if errExt := storage.ValidateExtension(filename); errExt != nil {
		return nil, errs.NewBadRequestError(errExt.Error())
	}

	// 3. Validate MIME Type
	mimeType, errMIME := storage.ValidateMIMEType(content)
	if errMIME != nil {
		return nil, errs.NewBadRequestError(errMIME.Error())
	}

	// 4. Calculate SHA-256 checksum
	checksum := storage.CalculateSHA256(content)

	// 5. Determine version number for document type
	maxVersion, errVer := s.docRepo.GetMaxVersionForDocumentType(ctx, orgID, docType)
	if errVer != nil {
		maxVersion = 0
	}
	nextVersion := maxVersion + 1

	// 6. Generate Document ID and Storage Key
	docID := uuid.New()
	storageKey := storage.BuildDeterministicStorageKey(orgID.String(), docType, docID.String(), nextVersion)

	callerUUID, _ := uuid.Parse(callerID)

	doc := &domain.OrganizationDocument{
		ID:               docID,
		OrganizationID:   orgID,
		DocumentType:     strings.ToLower(strings.TrimSpace(docType)),
		OriginalFilename: filename,
		StorageKey:       storageKey,
		MIMEType:         mimeType,
		FileSizeBytes:    int64(len(content)),
		ChecksumSHA256:   checksum,
		UploadedBy:       callerUUID,
		Status:           domain.DocumentStatusPending,
		Version:          nextVersion,
	}

	// 7. Upload Binary Object to Storage Service
	errStore := s.storageService.PutObject(ctx, storageKey, storage.NewMemoryBufferReader(content), int64(len(content)), mimeType)
	if errStore != nil {
		return nil, fmt.Errorf("failed to upload object to storage: %w", errStore)
	}

	// 8. Persist Document Metadata in PostgreSQL
	errDB := s.docRepo.CreateDocument(ctx, doc)
	if errDB != nil {
		// Non-ACID Cleanup: If database persistence fails after object upload, attempt cleanup of orphaned object
		_ = s.storageService.DeleteObject(ctx, storageKey)
		s.server.Logger.Error().Err(errDB).Str("storage_key", storageKey).Msg("failed to insert document metadata; orphaned storage object cleaned up")
		return nil, fmt.Errorf("failed to persist document metadata: %w", errDB)
	}

	// 9. Write Security Audit Log Event
	stmtAudit := `
		INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
		VALUES ($1, 'ORGANIZATION_DOCUMENT_UPLOADED', 'ORGANIZATION_DOCUMENT', $2, $3, $4, jsonb_build_object('document_type', $5, 'filename', $6, 'version', $7, 'checksum', $8))
	`
	_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), docID.String(), callerID, orgID.String(), doc.DocumentType, filename, nextVersion, checksum)

	return doc, nil
}

func (s *OrganizationDocumentApplicationService) ListDocuments(ctx context.Context, callerID string, orgID uuid.UUID) ([]DocumentWithPresignedURL, error) {
	if callerID == "" {
		return nil, errs.NewUnauthorizedError("authentication required")
	}

	// Verify caller membership or platform staff role
	var isMember bool
	_ = s.server.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships 
			WHERE user_id = $1::uuid AND organization_id = $2::uuid AND is_active = TRUE
		)
	`, callerID, orgID.String()).Scan(&isMember)

	if !isMember {
		var isPAdmin bool
		_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, callerID).Scan(&isPAdmin)
		if !isPAdmin {
			return nil, errs.NewForbiddenError("Access denied to organization documents")
		}
	}

	docs, err := s.docRepo.ListDocumentsByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var result []DocumentWithPresignedURL
	for _, d := range docs {
		url, errURL := s.storageService.GeneratePresignedURL(ctx, d.StorageKey, 15*time.Minute)
		if errURL != nil {
			url = ""
		}
		result = append(result, DocumentWithPresignedURL{
			Document:     d,
			PresignedURL: url,
		})
	}
	if result == nil {
		result = []DocumentWithPresignedURL{}
	}

	return result, nil
}

func (s *OrganizationDocumentApplicationService) ReviewDocument(
	ctx context.Context,
	reviewerID string,
	docID uuid.UUID,
	status domain.OrganizationDocumentStatus,
	rejectionReason *string,
) error {
	if reviewerID == "" {
		return errs.NewUnauthorizedError("authentication required")
	}

	// 1. Verify caller is authorized Platform Staff
	var isPAdmin bool
	_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, reviewerID).Scan(&isPAdmin)
	if !isPAdmin {
		return errs.NewForbiddenError("Only authorized platform staff may review organization documents")
	}

	// 2. Fetch document
	doc, errDoc := s.docRepo.GetDocumentByID(ctx, docID)
	if errDoc != nil {
		return errs.NewNotFoundError("document not found")
	}

	// 3. Prevent Self-Review (Caller cannot be the same user who uploaded the document)
	reviewerUUID, _ := uuid.Parse(reviewerID)
	if doc.UploadedBy == reviewerUUID {
		return errs.NewForbiddenError("Platform staff cannot review their own uploaded organization documents")
	}

	if status != domain.DocumentStatusApproved && status != domain.DocumentStatusRejected {
		return errs.NewBadRequestError("Review status must be either 'approved' or 'rejected'")
	}

	if status == domain.DocumentStatusRejected && (rejectionReason == nil || strings.TrimSpace(*rejectionReason) == "") {
		return errs.NewBadRequestError("Rejection reason is required when rejecting a document")
	}

	// 4. Update Document Review in PostgreSQL
	errUpdate := s.docRepo.UpdateDocumentReview(ctx, docID, status, reviewerUUID, rejectionReason)
	if errUpdate != nil {
		return fmt.Errorf("failed to update document review status: %w", errUpdate)
	}

	// 5. Log Audit Event
	reasonStr := ""
	if rejectionReason != nil {
		reasonStr = *rejectionReason
	}
	stmtAudit := `
		INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
		VALUES ($1, 'ORGANIZATION_DOCUMENT_REVIEWED', 'ORGANIZATION_DOCUMENT', $2, $3, $4, jsonb_build_object('status', $5, 'reason', $6))
	`
	_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), docID.String(), reviewerID, doc.OrganizationID.String(), string(status), reasonStr)

	return nil
}

func (s *OrganizationDocumentApplicationService) ApproveOrganization(ctx context.Context, reviewerID string, orgID uuid.UUID) error {
	if reviewerID == "" {
		return errs.NewUnauthorizedError("authentication required")
	}

	// 1. Verify Platform Staff Authority
	var isPAdmin bool
	_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, reviewerID).Scan(&isPAdmin)
	if !isPAdmin {
		return errs.NewForbiddenError("Only authorized platform staff may approve healthcare organization activation")
	}

	// 2. Fetch Organization to verify it exists
	org, errOrg := s.orgRepo.GetByID(ctx, orgID)
	if errOrg != nil {
		return errs.NewNotFoundError("organization not found")
	}

	// 3. Verify Required Documents for Organization Category are Approved
	approvedTypes, errTypes := s.docRepo.GetApprovedDocumentTypes(ctx, orgID)
	if errTypes != nil {
		return fmt.Errorf("failed to verify approved document types: %w", errTypes)
	}

	if len(approvedTypes) == 0 {
		return errs.NewBadRequestError("Organization activation rejected: At least one approved regulatory document (e.g. registration_certificate, operating_license, or medical_license) is required before activation.")
	}

	// 4. Update Organization Status to 'active'
	_, _ = s.server.DB.Pool.Exec(ctx, `UPDATE organization.organizations SET status = 'active', updated_at = CURRENT_TIMESTAMP WHERE id = $1::uuid`, orgID.String())

	// 5. Emit Audit Event
	stmtAudit := `
		INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
		VALUES ($1, 'ORGANIZATION_VERIFICATION_APPROVED', 'ORGANIZATION', $2, $3, $2, jsonb_build_object('previous_status', $4, 'new_status', 'active'))
	`
	_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), orgID.String(), reviewerID, org.Status)

	return nil
}

func (s *OrganizationDocumentApplicationService) RejectOrganization(ctx context.Context, reviewerID string, orgID uuid.UUID, reason string) error {
	if reviewerID == "" {
		return errs.NewUnauthorizedError("authentication required")
	}
	if strings.TrimSpace(reason) == "" {
		return errs.NewBadRequestError("Rejection reason is required")
	}

	var isPAdmin bool
	_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, reviewerID).Scan(&isPAdmin)
	if !isPAdmin {
		return errs.NewForbiddenError("Only authorized platform staff may reject healthcare organization activation")
	}

	_, err := s.server.DB.Pool.Exec(ctx, `UPDATE organization.organizations SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE id = $1::uuid`, orgID.String())
	if err != nil {
		return fmt.Errorf("failed to update organization status to rejected: %w", err)
	}

	stmtAudit := `
		INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
		VALUES ($1, 'ORGANIZATION_VERIFICATION_REJECTED', 'ORGANIZATION', $2, $3, $2, jsonb_build_object('reason', $4))
	`
	_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), orgID.String(), reviewerID, reason)

	return nil
}

// DownloadDocumentStream retrieves object binary stream for downloading.
func (s *OrganizationDocumentApplicationService) DownloadDocumentStream(ctx context.Context, callerID string, docID uuid.UUID) (domain.OrganizationDocument, []byte, error) {
	doc, err := s.docRepo.GetDocumentByID(ctx, docID)
	if err != nil {
		return domain.OrganizationDocument{}, nil, errs.NewNotFoundError("document not found")
	}

	// Verify membership or platform admin
	var isMember bool
	_ = s.server.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships 
			WHERE user_id = $1::uuid AND organization_id = $2::uuid AND is_active = TRUE
		)
	`, callerID, doc.OrganizationID.String()).Scan(&isMember)

	if !isMember {
		var isPAdmin bool
		_ = s.server.DB.Pool.QueryRow(ctx, `SELECT COALESCE(is_platform_admin, FALSE) FROM identity.users WHERE id = $1::uuid`, callerID).Scan(&isPAdmin)
		if !isPAdmin {
			return domain.OrganizationDocument{}, nil, errs.NewForbiddenError("Access denied to document")
		}
	}

	rc, errGet := s.storageService.GetObject(ctx, doc.StorageKey)
	if errGet != nil {
		return domain.OrganizationDocument{}, nil, fmt.Errorf("failed to read storage object: %w", errGet)
	}
	defer rc.Close()

	data, errRead := io.ReadAll(rc)
	if errRead != nil {
		return domain.OrganizationDocument{}, nil, fmt.Errorf("failed to read document object stream: %w", errRead)
	}

	return *doc, data, nil
}
