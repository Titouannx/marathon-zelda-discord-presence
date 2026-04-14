//go:build windows

package platform

import (
	"errors"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
const runValueName = "MarathonZeldaDiscordPresence"
const uninstallKeyPath = `Software\Microsoft\Windows\CurrentVersion\Uninstall\MarathonZeldaDiscordPresence`

func IsAutoStartEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	value, _, err := key.GetStringValue(runValueName)
	if err != nil {
		return false
	}

	executable, err := executablePath()
	if err != nil {
		return false
	}

	return value == executable
}

func InstallAutoStart() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	executable, err := executablePath()
	if err != nil {
		return err
	}

	return key.SetStringValue(runValueName, executable)
}

func RemoveAutoStart() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return err
	}
	defer key.Close()

	if err := key.DeleteValue(runValueName); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

func EnsureAppRegistration(displayName string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, uninstallKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	executable, err := executablePath()
	if err != nil {
		return err
	}

	if err := key.SetStringValue("DisplayName", displayName); err != nil {
		return err
	}
	if err := key.SetStringValue("Publisher", "Titouannx"); err != nil {
		return err
	}
	if err := key.SetStringValue("DisplayIcon", executable); err != nil {
		return err
	}
	if err := key.SetStringValue("InstallLocation", filepath.Dir(executable)); err != nil {
		return err
	}

	uninstallCommand := `"` + executable + `" --uninstall`
	if err := key.SetStringValue("UninstallString", uninstallCommand); err != nil {
		return err
	}
	if err := key.SetStringValue("QuietUninstallString", uninstallCommand); err != nil {
		return err
	}
	if err := key.SetDWordValue("NoModify", 1); err != nil {
		return err
	}
	if err := key.SetDWordValue("NoRepair", 1); err != nil {
		return err
	}

	return nil
}

func RemoveAppRegistration() error {
	if err := registry.DeleteKey(registry.CURRENT_USER, uninstallKeyPath); err != nil &&
		!errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}
