//go:build darwin

package platform

import "os/exec"

func ShowInfo(title string, message string) error {
	return exec.Command("osascript", "-e", `display dialog "`+escapeAppleScript(message)+`" with title "`+escapeAppleScript(title)+`" buttons {"OK"} default button "OK"`).Run()
}

func ShowError(title string, message string) error {
	return exec.Command("osascript", "-e", `display dialog "`+escapeAppleScript(message)+`" with title "`+escapeAppleScript(title)+`" buttons {"OK"} default button "OK" with icon caution`).Run()
}

func escapeAppleScript(value string) string {
	replacer := map[rune]string{
		'\\': `\\`,
		'"':  `\"`,
	}
	result := make([]rune, 0, len(value))
	for _, r := range value {
		if replacement, ok := replacer[r]; ok {
			for _, rr := range replacement {
				result = append(result, rr)
			}
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
