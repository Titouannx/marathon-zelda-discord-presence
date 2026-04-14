//go:build windows

package assets

import (
	"encoding/binary"
	"testing"
)

type icoEntry struct {
	Width  int
	Height int
}

func TestTrayIconIsWindowsIco(t *testing.T) {
	if len(TrayIcon) < 6 {
		t.Fatalf("TrayIcon should not be empty")
	}

	if TrayIcon[0] != 0x00 || TrayIcon[1] != 0x00 || TrayIcon[2] != 0x01 || TrayIcon[3] != 0x00 {
		t.Fatalf("TrayIcon should embed a valid ICO header on Windows")
	}
}

func TestTrayIconEmbedsMultipleSizes(t *testing.T) {
	count := int(binary.LittleEndian.Uint16(TrayIcon[4:6]))
	if count < 4 {
		t.Fatalf("expected a multi-size ICO with at least 4 entries, got %d", count)
	}

	entries := parseIcoEntries(t, TrayIcon)
	if !hasIcoSize(entries, 16, 16) {
		t.Fatalf("expected ICO to embed a 16x16 image, got %#v", entries)
	}
	if !hasIcoSize(entries, 32, 32) {
		t.Fatalf("expected ICO to embed a 32x32 image, got %#v", entries)
	}
}

func parseIcoEntries(t *testing.T, payload []byte) []icoEntry {
	t.Helper()

	count := int(binary.LittleEndian.Uint16(payload[4:6]))
	entries := make([]icoEntry, 0, count)
	for i := 0; i < count; i++ {
		offset := 6 + (i * 16)
		if len(payload) < offset+16 {
			t.Fatalf("ICO payload truncated at entry %d", i)
		}

		width := int(payload[offset])
		if width == 0 {
			width = 256
		}
		height := int(payload[offset+1])
		if height == 0 {
			height = 256
		}
		entries = append(entries, icoEntry{Width: width, Height: height})
	}

	return entries
}

func hasIcoSize(entries []icoEntry, width int, height int) bool {
	for _, entry := range entries {
		if entry.Width == width && entry.Height == height {
			return true
		}
	}
	return false
}
