//go:build darwin

package platform

import (
	"os"
	"os/user"
	"path/filepath"
)

const launchAgentName = "bzh.loon.marathon-zelda-discord-presence.plist"

func launchAgentPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(currentUser.HomeDir, "Library", "LaunchAgents", launchAgentName), nil
}

func IsAutoStartEnabled() bool {
	path, err := launchAgentPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

func InstallAutoStart() error {
	executable, err := executablePath()
	if err != nil {
		return err
	}

	path, err := launchAgentPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>bzh.loon.marathon-zelda-discord-presence</string>
  <key>ProgramArguments</key>
  <array>
    <string>` + executable + `</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <false/>
</dict>
</plist>`

	return os.WriteFile(path, []byte(plist), 0o644)
}

func RemoveAutoStart() error {
	path, err := launchAgentPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
