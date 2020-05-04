package infras

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	// "regexp"
	// "strings"
	"io/ioutil"
	"path/filepath"
	"sync"

	// "github.com/astaxie/beego/logs"

	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/utils"
)

var fileMutex = new(sync.Mutex)

func ListDir(rel string) []*models.File {
	ch := make(chan []*models.File)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(config.GetSharedDir(), rel)

		ii, err := ioutil.ReadDir(target)
		if err != nil {
			panic(fmt.Sprintf("Can not list dir %s: %s", target, err.Error()))
		}

		files := make([]*models.File, 0)
		for _, i := range ii {
			file := models.NewFile(i.Name(), i.Size(), i.ModTime(), i.IsDir())
			files = append(files, file)
		}
		ch <- files
	}()
	return <-ch
}

func GetStat(rel string) os.FileInfo {
	ch := make(chan os.FileInfo)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(config.GetSharedDir(), rel)
		stat, err := os.Stat(target)
		if err != nil {
			ch <- nil
			return
		}
		ch <- stat
	}()
	return <-ch
}

func DeleteFile(rel string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(config.GetSharedDir(), rel)
		err := os.Remove(target)
		if err != nil {
			ch <- fmt.Errorf("Could not remove item: %s", err.Error())
			return
		}
		ch <- nil
	}()
	return <-ch
}

func SaveToFile(file multipart.File, rel string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(config.GetSharedDir(), rel)
		dst, err := os.Create(target)
		defer dst.Close()
		if err != nil {
			ch <- err
			return
		}
		_, err = io.Copy(dst, file)
		if err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()
	return <-ch
}

func Mkdir(rel string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		target := filepath.Join(config.GetSharedDir(), rel)
		stat, err := os.Stat(target)
		if err == nil {
			var w string
			if stat.IsDir() {
				w = "directory"
			} else {
				w = "file"
			}
			ch <- fmt.Errorf("A %s already exists in %s", w, rel)
			return
		}
		err = utils.EnsureDir(target)
		if err != nil {
			ch <- fmt.Errorf("Failed to maked directory %s: %s", rel, err.Error())
			return
		}
		ch <- nil
	}()
	return <-ch
}

func RenameFile(src, dest string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		sharedDir, err := filepath.Abs(config.GetSharedDir())
		if err != nil {
			ch <- err
			return
		}
		srcPath := filepath.Join(sharedDir, src)
		destPath := filepath.Join(sharedDir, dest)
		if !utils.FileExists(srcPath) {
			ch <- fmt.Errorf("File does not exist in path `%s`", srcPath)
			return
		}
		if utils.FileExists(destPath) {
			ch <- fmt.Errorf("File already exists in path `%s`", destPath)
			return
		}

		destBase, err := filepath.Abs(filepath.Dir(destPath))
		if err != nil {
			ch <- err
			return
		}
		stat, err := os.Stat(destBase)
		if !(err == nil && stat.IsDir()) {
			ch <- fmt.Errorf("Dest dir(%s) is not exists", destBase)
			return
		}

		under := strings.HasPrefix(destBase, sharedDir)
		if !under {
			ch <- fmt.Errorf("Tried to save file under shared dir.")
			return
		}

		err = os.Rename(srcPath, destPath)
		if err != nil {
			ch <- fmt.Errorf("Failed to rename: %s -> %s", srcPath, destPath)
			return
		}
		ch <- nil
	}()
	return <-ch
}
