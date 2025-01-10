package shared

import (
	"errors"
	"os"
)

func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return info.IsDir()
}