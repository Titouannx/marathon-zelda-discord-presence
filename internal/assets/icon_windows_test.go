//go:build windows

package assets

import "testing"

func TestTrayIconIsWindowsIco(t *testing.T) {
	if len(TrayIcon) < 4 {
		t.Fatalf("TrayIcon should not be empty")
	}

	if TrayIcon[0] != 0x00 || TrayIcon[1] != 0x00 || TrayIcon[2] != 0x01 || TrayIcon[3] != 0x00 {
		t.Fatalf("TrayIcon should embed a valid ICO header on Windows")
	}
}
