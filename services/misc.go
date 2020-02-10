package services

import (
	"github.com/endaaman/api.endaaman.me/infras"
)

func RetrieveWarnings(ch chan<- []string) {
	infras.WaitReader()
    ch <- infras.GetCachedWarnings()
}
