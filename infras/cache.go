package infras

import (
	"sync"
	"github.com/endaaman/api.endaaman.me/models"
)

var aa_mutex sync.RWMutex
var ww_mutex sync.RWMutex
var loading = true
var articles []*models.Article
var warnings []string

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

// func AddCachedArticle(a *models.Article) {
// 	aa := GetCachedArticles()
// 	aa = append(aa, a)
// 	SetCachedArticles(aa)
// }

func GetCachedWarnings() []string {
	ww_mutex.RLock()
	var ww = warnings
	ww_mutex.RUnlock()
	return ww
}

func SetCachedWarnings(ww []string) {
	ww_mutex.Lock()
	warnings = ww
	ww_mutex.Unlock()
}

// func AddCachedWarning(w string) {
// 	ww := GetCachedWarnings()
// 	ww = append(ww, w)
// 	SetCachedWarnings(ww)
// }
