//go:build !windows

package assets

import _ "embed"

//go:embed tray_icon_nonwindows.png
var TrayIcon []byte
