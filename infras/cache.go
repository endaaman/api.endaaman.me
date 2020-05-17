package infras

import (
	"sync"

	"github.com/endaaman/api.endaaman.me/models"
)

type lockable struct {
	mutex *sync.RWMutex
	value interface{}
}

func newLockcable(value interface{}) *lockable {
	l := lockable{}
	l.mutex = new(sync.RWMutex)
	l.value = value
	return &l
}

func (l *lockable) get() interface{} {
	l.mutex.RLock()
	var v = l.value
	l.mutex.RUnlock()
	return v
}

func (l *lockable) set(v interface{}) {
	l.mutex.Lock()
	l.value = v
	l.mutex.Unlock()
}

var aa = newLockcable(make([]*models.Article, 0))
var cc = newLockcable(make([]*models.Category, 0))
var ww = newLockcable(make(map[string][]string))

func GetCachedArticles() []*models.Article {
	v, ok := aa.get().([]*models.Article)
	if !ok {
		panic("Invalid type")
	}
	return v
}

func SetCachedArticles(v []*models.Article) {
	aa.set(v)
}

func GetCachedCategorys() []*models.Category {
	v, ok := cc.get().([]*models.Category)
	if !ok {
		panic("Invalid type")
	}
	return v
}

func SetCachedCategorys(v []*models.Category) {
	cc.set(v)
}

func GetCachedWarnings() map[string][]string {
	v, ok := ww.get().(map[string][]string)
	if !ok {
		panic("Invalid type")
	}
	return v
}

func SetCachedWarnings(v map[string][]string) {
	ww.set(v)
}
