package infras

import (
	"sync"
	"github.com/endaaman/api.endaaman.me/models"
)

type lockable struct {
	mutex *sync.Mutex
}

var aa_mutex sync.RWMutex
var cc_mutex sync.RWMutex
var loading = true
var articles []*models.Article
var categorys []*models.Category

func GetCachedArticles() []*models.Article {
	aa_mutex.RLock()
	var aa = articles
	aa_mutex.RUnlock()
	return aa
}

func SetCachedArticles(aa []*models.Article) {
	aa_mutex.Lock()
	articles = aa
	aa_mutex.Unlock()
}

func GetCachedCategorys() []*models.Category {
	cc_mutex.RLock()
	var cc = categorys
	cc_mutex.RUnlock()
	return cc
}

func SetCachedCategorys(cc []*models.Category) {
	cc_mutex.Lock()
	categorys = cc
	cc_mutex.Unlock()
}
