package job

type Task struct {
	TypeName string
	Payload  map[string]interface{}
}

func NewVerifyEmailTask(email, name, link, _, _ string) (*Task, error) {
	return &Task{
		TypeName: "verify_email",
		Payload: map[string]interface{}{
			"email": email,
			"name":  name,
			"link":  link,
		},
	}, nil
}

func NewLoginOTPEmailTask(email, name, otp string) (*Task, error) {
	return &Task{
		TypeName: "login_otp",
		Payload: map[string]interface{}{
			"email": email,
			"name":  name,
			"otp":   otp,
		},
	}, nil
}

func NewForgotPasswordEmailTask(email, link string) (*Task, error) {
	return &Task{
		TypeName: "forgot_password",
		Payload: map[string]interface{}{
			"email": email,
			"link":  link,
		},
	}, nil
}

func NewInviteMemberEmailTask(email, inviterName, tenantName, link string) (*Task, error) {
	return &Task{
		TypeName: "invite_member",
		Payload: map[string]interface{}{
			"email":       email,
			"inviterName": inviterName,
			"tenantName":  tenantName,
			"link":        link,
		},
	}, nil
}

type AuditLogTaskPayload struct {
	IsPlatform   bool
	TenantID     *string
	ActorID      *string
	ActorRole    *string
	Action       string
	ResourceType *string
	ResourceID   *string
	ResourceName *string
	Details      string
	Severity     string
	Status       string
	Category     string
	UserID       *string
	IPAddress    *string
	UserAgent    *string
	RequestID    *string
	BeforeState  *string
	AfterState   *string
}

func NewAuditLogTask(_ interface{}) (*Task, error) {
	return &Task{TypeName: "audit_log"}, nil
}

func NewSignatureCreatedTask(sigID, userID, tenantID, imageUrl string) (*Task, error) {
	return &Task{
		TypeName: "signature_created",
		Payload: map[string]interface{}{
			"id":                sigID,
			"userId":            userID,
			"tenantId":          tenantID,
			"signatureImageUrl": imageUrl,
		},
	}, nil
}
