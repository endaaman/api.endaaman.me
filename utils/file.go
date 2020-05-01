package utils

import (
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}
