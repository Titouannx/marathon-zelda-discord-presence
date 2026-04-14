//go:build !windows && !darwin

package platform

func ShowInfo(title string, message string) error {
	return nil
}

func ShowError(title string, message string) error {
	return nil
}
