//go:build windows

package assets

import _ "embed"

//go:embed tray_icon_windows.ico
var TrayIcon []byte
