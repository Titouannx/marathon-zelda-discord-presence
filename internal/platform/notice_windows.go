//go:build windows

package platform

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL       = windows.NewLazySystemDLL("user32.dll")
	messageBoxWProc = user32DLL.NewProc("MessageBoxW")
)

const (
	messageBoxOK              = 0x00000000
	messageBoxIconError       = 0x00000010
	messageBoxIconInformation = 0x00000040
	messageBoxTopMost         = 0x00040000
)

func ShowInfo(title string, message string) error {
	return showMessageBox(title, message, messageBoxOK|messageBoxIconInformation|messageBoxTopMost)
}

func ShowError(title string, message string) error {
	return showMessageBox(title, message, messageBoxOK|messageBoxIconError|messageBoxTopMost)
}

func showMessageBox(title string, message string, flags uintptr) error {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	messagePtr, err := windows.UTF16PtrFromString(message)
	if err != nil {
		return err
	}

	_, _, callErr := messageBoxWProc.Call(
		0,
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		flags,
	)
	if callErr != windows.ERROR_SUCCESS && callErr != nil {
		return callErr
	}
	return nil
}
