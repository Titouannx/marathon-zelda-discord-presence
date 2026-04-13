package platform

import "os"

func executablePath() (string, error) {
	return os.Executable()
}
