package infras

import (
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/bep/debounce"
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/utils"
	"github.com/radovskyb/watcher"
)

var LastError error = nil
var IsWatcherActive bool = false
var mutex sync.Mutex
var ch = make(chan bool)

func notify() {
	logs.Info("Detect changes")
	ReadAllArticles()
	WaitIO()
	ch <- true
}

func AwaitNextChange() {
	logs.Info("Start awaiting next change and loading done")
	select {
	case <-ch:
		logs.Info("Load done by event triggered")
	case <-time.After(3 * time.Second):
	}
}

func StartWatcher() {
	if IsWatcherActive {
		logs.Warn("Tried to start watcher twice")
		return
	}
	mutex.Lock()
	LastError = nil
	IsWatcherActive = true
	ch := make(chan error)
	go watch(ch)
	LastError = <-ch
	IsWatcherActive = false
	logs.Info("Watcher closed")
	mutex.Unlock()
}

func watch(ch chan<- error) {
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Rename, watcher.Move, watcher.Write)

	go func() {
		notify() // run as first
		debounced := debounce.New(time.Millisecond * 100)
		for {
			select {
			case <-w.Event:
				debounced(notify)
			case err := <-w.Error:
				ch <- fmt.Errorf("Error occured on watcher: %s", err.Error())
				w.Close()
				return
			case <-w.Closed:
				logs.Info("Watcher has been closed")
				return
			}
		}
	}()

	articlesDir := config.GetArticlesDir()
	if !utils.IsDir(articlesDir) {
		ch <- fmt.Errorf("articles dir(%s) is not directory", articlesDir)
		return
	}

	err := w.AddRecursive(articlesDir)
	if err != nil {
		ch <- fmt.Errorf("Failed to add recursive watching: %s", err.Error())
		return
	}

	err = w.Start(time.Millisecond * 300)
	if err != nil {
		ch <- fmt.Errorf("Watcher has been closed with error: %s", err.Error())
		return
	}
	logs.Info("Watcher has closed")
	ch <- nil
}
