//go:build windows

package platform

import "golang.org/x/sys/windows/registry"

const runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
const runValueName = "MarathonZeldaDiscordPresence"

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
		return err
	}
	defer key.Close()

	return key.DeleteValue(runValueName)
}
