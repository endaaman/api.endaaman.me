package infras

import (
	"fmt"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/bep/debounce"
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/utils"
	"github.com/radovskyb/watcher"
)

var LastError error = nil
var watcherInstance *watcher.Watcher = nil
var watcherMutex = &sync.Mutex{}
var changeNotify chan bool
var closeNotify chan error

func notify() {
	logs.Info("Detect changes")
	ReadAllArticles()
	WaitIO()
	if changeNotify != nil {
		changeNotify <- true
	}
}

func AwaitNextChange() {
	logs.Info("Start awaiting next change and loading done")
	changeNotify = make(chan bool)
	select {
	case <-changeNotify:
		logs.Info("Load done by event triggered")
	case <-time.After(3 * time.Second):
		logs.Warn("Timeout notify")
	}
}

func IsWatcherActive() bool {
	return watcherInstance != nil
}

func StartWatcher() {
	watcherMutex.Lock()
	logs.Info("Starting watcher")

	LastError = nil
	closeNotify = make(chan error)
	watcherInstance = watcher.New()
	watcherInstance.FilterOps(watcher.Create, watcher.Rename, watcher.Move, watcher.Write)

	go func() {
		articlesDir := config.GetArticlesDir()
		if !utils.IsDir(articlesDir) {
			closeNotify <- fmt.Errorf("articles dir(%s) is not directory", articlesDir)
			return
		}

		err := watcherInstance.AddRecursive(articlesDir)
		if err != nil {
			closeNotify <- fmt.Errorf("Failed to add recursive watching: %s", err.Error())
			return
		}

		go func() {
			logs.Info("Watcher event loop started")
			notify() // run at first
			debounced := debounce.New(time.Millisecond * 100)
			for {
				select {
				case <-watcherInstance.Event:
					debounced(notify)
				case err := <-watcherInstance.Error:
					closeNotify <- fmt.Errorf("Error occured on watcher: %s", err.Error())
					return
					// case <-watcherInstance.Closed:
					// 	closeNotify <- nil
					// 	return
				}
			}
		}()

		err = watcherInstance.Start(time.Millisecond * 300)
		if err != nil {
			closeNotify <- fmt.Errorf("Watcher has been closed with error: %s", err.Error())
			return
		}
		closeNotify <- nil
	}()
	LastError = <-closeNotify
	if LastError != nil {
		logs.Error("Watcher closed because: %s", LastError.Error())
	} else {
		logs.Warn("Watcher closed without any error")
	}
	watcherInstance.Close()
	watcherInstance = nil

	logs.Info("Watcher has successfull closed")
	watcherMutex.Unlock()
}

func RestartWatcher() {
	if IsWatcherActive() {
		logs.Info("Tried to close watcher")
		closeNotify <- fmt.Errorf("Close watcher manually")
	}
	StartWatcher()
}
