package infras

import (
    "os"
    // "regexp"
    "fmt"
    // "strings"
    "sync"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego"
    // "github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/models"
)

var fileMutex = new(sync.Mutex)

func ListDir(rel string, ch chan []*models.File) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	baseDir := beego.AppConfig.String("private_dir")
	target := filepath.Join(baseDir, rel)

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
}

func IsDir(rel string, ch chan bool) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	baseDir := beego.AppConfig.String("private_dir")
	target := filepath.Join(baseDir, rel)
	stat, err := os.Stat(target)
	if err != nil {
		ch<- false
		return
	}
	ch<- stat.IsDir()
}
