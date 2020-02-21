package services

import (
	"github.com/endaaman/api.endaaman.me/infras"
)

func GetWarnings() map[string][]string {
    return infras.GetCachedWarnings()
}
