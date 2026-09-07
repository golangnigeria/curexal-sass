package domain

import "errors"

var (
	ErrUnauthorizedBranchAccess = errors.New("unauthorized facility branch access")
	ErrUnassignedFacilityBranch = errors.New("unassigned facility branch: user has no operational branch assignment")
	ErrBranchSelectionRequired  = errors.New("active branch selection is required")
	ErrInvalidSelectionToken    = errors.New("invalid or expired branch selection token")
)
