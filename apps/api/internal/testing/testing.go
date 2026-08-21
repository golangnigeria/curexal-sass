package testing

// Testing Package: Shared test utilities, mocks, and test fixtures for Curexal V2.
type TestFixture struct {
	TenantSlug string
	UserEmail  string
}

func NewTestFixture(tenantSlug, userEmail string) *TestFixture {
	return &TestFixture{
		TenantSlug: tenantSlug,
		UserEmail:  userEmail,
	}
}
