package infras

import (
	"fmt"
	"log"
	"time"
	"github.com/radovskyb/watcher"
	"github.com/bep/debounce"
)


func notify() {
	fmt.Println("DEBOUCED")
}

func StartWatching() {
	w := watcher.New()

	// w.SetMaxEvents(1)
	w.FilterOps(watcher.Create, watcher.Rename, watcher.Move, watcher.Write)
	// r := regexp.MustCompile("^abc$")
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		debounced := debounce.New(time.Second)
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				debounced(notify)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.AddRecursive("./shared/articles"); err != nil {
		log.Fatalln(err)
	}

	// for path, f := range w.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
