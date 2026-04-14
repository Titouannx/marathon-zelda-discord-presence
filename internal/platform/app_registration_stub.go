//go:build !windows

package platform

func EnsureAppRegistration(_ string) error {
	return nil
}

func RemoveAppRegistration() error {
	return nil
}
