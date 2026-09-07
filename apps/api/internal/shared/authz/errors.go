package authz

import "errors"

var (
	ErrUnauthorized       = errors.New("unauthorized: missing or invalid credentials")
	ErrForbidden          = errors.New("forbidden: insufficient permissions for requested resource")
	ErrScopeMismatch      = errors.New("forbidden: request scope does not match principal context")
	ErrProductDisabled    = errors.New("forbidden: product is not enabled for this facility")
	ErrBranchNotAssigned  = errors.New("forbidden: user is not assigned to this branch")
	ErrCredentialRequired = errors.New("forbidden: professional credential required for this clinical action")
)
