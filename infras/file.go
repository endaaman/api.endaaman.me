package infras

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/utils"
)

var fileMutex = new(sync.Mutex)

func prepareRelative(rel string) (sharedDir string, target string, err error) {
	sharedDir = config.GetSharedDir()
	target = filepath.Join(sharedDir, rel)
	if !utils.IsDir(sharedDir) {
		err = fmt.Errorf("Shared dir(%s) does not exist.", sharedDir)
		return
	}
	if !utils.IsUnder(sharedDir, target) {
		err = fmt.Errorf("Tried to access under shared dir.")
		return
	}
	return
}

func ListDir(rel string) (files []*models.File, err error) {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, target, err := prepareRelative(rel)
		if err != nil {
			ch <- true
			return
		}

		items, err := ioutil.ReadDir(target)
		if err != nil {
			err = fmt.Errorf("Can not list dir %s: %s", target, err.Error())
			ch <- true
		}

		files = make([]*models.File, 0)
		for _, i := range items {
			file := models.NewFile(i.Name(), i.Size(), i.ModTime(), i.IsDir())
			files = append(files, file)
		}
		ch <- true
	}()
	<-ch
	return
}

func getStat(rel string) (stat os.FileInfo, err error) {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, target, err := prepareRelative(rel)
		if err != nil {
			ch <- true
			return
		}

		stat, err = os.Stat(target)
		if err != nil {
			ch <- true
			return
		}
		ch <- true
	}()
	<-ch
	return
}

func FileExists(rel string) (bool, error) {
	stat, err := getStat(rel)
	return stat != nil, err
}

func IsDir(rel string) (bool, error) {
	stat, err := getStat(rel)
	return stat != nil && stat.IsDir(), err
}

func DeleteFile(rel string) (err error) {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, target, err := prepareRelative(rel)
		if err != nil {
			ch <- true
			return
		}

		err = os.Remove(target)
		if err != nil {
			err = fmt.Errorf("Could not remove item: %s", err.Error())
			return
		}
		ch <- true
	}()
	<-ch
	return
}

func SaveToFile(file multipart.File, rel string) (err error) {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, target, err := prepareRelative(rel)
		if err != nil {
			ch <- true
			return
		}

		dst, err := os.Create(target)
		defer dst.Close()
		if err != nil {
			ch <- true
			return
		}
		_, err = io.Copy(dst, file)
		if err != nil {
			ch <- true
			return
		}
		ch <- true
	}()
	<-ch
	return
}

func Mkdir(rel string) (err error) {
	ch := make(chan bool)
	go func() {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, target, err := prepareRelative(rel)
		if err != nil {
			ch <- true
			return
		}
		stat, err := os.Stat(target)
		if stat != nil {
			var w string
			if stat.IsDir() {
				w = "directory"
			} else {
				w = "file"
			}
			err = fmt.Errorf("A %s already exists in %s", w, rel)
			ch <- true
			return
		}
		err = utils.EnsureDir(target)
		if err != nil {
			err = fmt.Errorf("Failed to maked directory %s: %s", rel, err.Error())
			ch <- true
			return
		}
		ch <- true
	}()
	<-ch
	return
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

		if !utils.IsUnder(destBase, sharedDir) {
			ch <- fmt.Errorf("Tried to access under shared dir.")
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
