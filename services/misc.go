package services

import (
	"fmt"

	"github.com/endaaman/api.endaaman.me/infras"
)

func GetWarnings() map[string][]string {
	return infras.GetCachedWarnings()
}

func IsWatcherActive() bool {
	return infras.IsWatcherActive
}

func GetWathcerLastError() string {
	if infras.LastError == nil {
		return ""
	}
	return infras.LastError.Error()
}

func RestartWatcher() error {
	if infras.IsWatcherActive {
		return fmt.Errorf("Watcher already started.")
	}
	go infras.StartWatcher()
	return nil
}
