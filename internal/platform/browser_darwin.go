//go:build darwin

package platform

import "os/exec"

func OpenURL(url string) error {
	return exec.Command("open", url).Start()
}
