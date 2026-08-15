package azt002

import "os"

// Should NOT be flagged: non-test files (the provider runtime and the acceptance test
// framework) legitimately read the credentials they authenticate with
func runtimeClientId() string {
	return os.Getenv("ARM_CLIENT_ID")
}

func runtimeClientSecret() string {
	return os.Getenv("ARM_CLIENT_SECRET")
}

func runtimeClientSecretAlt() string {
	return os.Getenv("ARM_CLIENT_SECRET_ALT")
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
