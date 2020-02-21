package infras

import (
    "os"
    // "io"
    "fmt"
    // "regexp"
    // "strings"
    "sync"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego"
    // "github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/models"
)

var fileMutex = new(sync.Mutex)

func ListDir(rel string) []*models.File {
	ch := make(chan []*models.File)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(beego.AppConfig.String("private_dir"), rel)

		ii, err := ioutil.ReadDir(target)
		if err != nil {
			panic(fmt.Sprintf("Can not list dir %s: %s", target, err.Error()))
		}

		files := make([]*models.File, 0)
		for _, i := range ii {
			file := models.NewFile(i.Name(), i.Size(), i.ModTime(), i.IsDir())
			files = append(files, file)
		}
		ch<- files
	}()
	return <-ch
}

func IsDir(rel string) bool {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(beego.AppConfig.String("private_dir"), rel)
		stat, err := os.Stat(target)
		if err != nil {
			ch<- false
			return
		}
		ch<- stat.IsDir()
	}()
	return <-ch
}

func Remove(rel string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(beego.AppConfig.String("private_dir"), rel)
		err := os.Remove(target)
		if err != nil {
			ch<- fmt.Errorf("Could not remove item: %s", err.Error())
			return
		}
		ch<- nil
	}()
	return <-ch
}
