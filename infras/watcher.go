package infras

import (
	"log"
	"time"
	"sync"
	"github.com/radovskyb/watcher"
	"github.com/bep/debounce"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

var watcher_mutex sync.Mutex
var ch = make(chan bool)

func notify() {
	logs.Info("Detect changes")
	ReadAllArticles().Wait()
	ch <- true
}

func AwaitNextChange() {
    select {
    case <-ch:
    case <-time.After(3 * time.Second):
    }
}

func StartWatching() {
	w := watcher.New()
	// w.SetMaxEvents(1)
	w.FilterOps(watcher.Create, watcher.Rename, watcher.Move, watcher.Write)
	// r := regexp.MustCompile("^abc$")
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		debounced := debounce.New(time.Millisecond * 100)
		for {
			select {
			case <-w.Event:
				debounced(notify)
			case err := <-w.Error:
				logs.Error(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	dir := beego.AppConfig.String("articles_dir")
	if err := w.AddRecursive(dir); err != nil {
		log.Fatalln(err)
	}

	// for path, f := range w.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }

	if err := w.Start(time.Millisecond * 300); err != nil {
		log.Fatalln(err)
	}
}
