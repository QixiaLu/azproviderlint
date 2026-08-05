package azt002

import "os"

// Should be flagged: reusing the test credentials
func badClientId() string {
	return os.Getenv("ARM_CLIENT_ID") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_ID"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

func badClientSecret() string {
	return os.Getenv("ARM_CLIENT_SECRET") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_SECRET"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

func badClientSecretAlt() string {
	return os.Getenv("ARM_CLIENT_SECRET_ALT") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_SECRET_ALT"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

// Should NOT be flagged: other environment variables
func goodSubscriptionId() string {
	return os.Getenv("ARM_SUBSCRIPTION_ID")
}

// Should NOT be flagged: not os.Getenv
type fakeOs struct{}

func (fakeOs) Getenv(k string) string { return "" }

func goodOtherGetenv() string {
	var f fakeOs
	return f.Getenv("ARM_CLIENT_ID")
}
