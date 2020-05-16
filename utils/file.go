package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func IsDir(path string) bool {
	stat, _ := os.Stat(path)
	if stat == nil {
		return false
	}
	return stat.IsDir()
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func IsUnder(base, target string) bool {
	return strings.HasPrefix(filepath.Clean(target), filepath.Clean(base))
}
