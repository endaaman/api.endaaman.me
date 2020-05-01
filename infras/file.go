package infras

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	// "regexp"
	// "strings"
	"io/ioutil"
	"path/filepath"
	"sync"

	// "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/logs"
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

func SaveToFile(rel string, file multipart.File) error {
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

func RenameFile(src, dest string) error {
	ch := make(chan error)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		srcPath := filepath.Join(config.GetSharedDir(), src)
		destPath := filepath.Join(config.GetSharedDir(), dest)
		if !utils.FileExists(srcPath) {
			ch <- fmt.Errorf("File does not exist in path `%s`", srcPath)
			return
		}
		if utils.FileExists(destPath) {
			ch <- fmt.Errorf("File already exists in path `%s`", destPath)
			return
		}
		// check if not deeper than root
		// check if dest base dir exists
		logs.Warn("RENAME %s -> %s", srcPath, destPath)
		ch <- nil
	}()
	return <-ch
}
