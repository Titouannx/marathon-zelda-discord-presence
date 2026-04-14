//go:build windows

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecutableEmbedsWindowsResourceObject(t *testing.T) {
	path := filepath.Join("resource_windows_amd64.syso")
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected Windows resource object at %s: %v", path, err)
	}
	if len(payload) < 2 {
		t.Fatalf("expected Windows resource object to be non-empty")
	}
	if payload[0] != 0x64 || payload[1] != 0x86 {
		t.Fatalf("expected AMD64 COFF resource object header, got %#x %#x", payload[0], payload[1])
	}
}
