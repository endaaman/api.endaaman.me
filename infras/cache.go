package infras

import (
	"sync"

	"github.com/endaaman/api.endaaman.me/models"
)

var cache_mutex sync.RWMutex
var loading = true
var articles []*models.Article

func GetCachedArticles() []*models.Article {
	cache_mutex.RLock()
	var aa = articles
	cache_mutex.RUnlock()
	return aa
}

func SetCachedArticles(aa []*models.Article) {
	cache_mutex.Lock()
	articles = aa
	cache_mutex.Unlock()
}

func IsLoading() bool {
	cache_mutex.RLock()
	var i = loading
	cache_mutex.RUnlock()
	return i
}
