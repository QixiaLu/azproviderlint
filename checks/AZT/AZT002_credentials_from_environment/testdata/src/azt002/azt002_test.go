package azt002

import "os"

// Should be flagged: reusing the test credentials in a test file
func badClientId() string {
	return os.Getenv("ARM_CLIENT_ID") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_ID"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

func badClientSecret() string {
	return os.Getenv("ARM_CLIENT_SECRET") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_SECRET"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

func badClientSecretAlt() string {
	return os.Getenv("ARM_CLIENT_SECRET_ALT") // want `tests should not obtain credentials via os\.Getenv\("ARM_CLIENT_SECRET_ALT"\), create an azurerm_user_assigned_identity with minimal permissions as part of the test configuration instead`
}

// Should NOT be flagged: other environment variables in a test file
func goodTestSubscriptionId() string {
	return os.Getenv("ARM_SUBSCRIPTION_ID")
}
