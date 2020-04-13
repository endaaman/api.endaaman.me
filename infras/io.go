package infras

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/utils"
)

const CATEGORY_FILE_NAME = "meta.json"
const (
	FILE_TYPE_ARTICLE = iota
	FILE_TYPE_CATEGORY
	FILE_TYPE_OTHER
)

var ioWaiter = new(sync.WaitGroup)
var ioMutex = new(sync.Mutex)

func WaitIO() {
	ioWaiter.Wait()
}

type fileItem struct {
	baseDir  string
	relative string
	file     os.FileInfo
}

func dirwalk(base, dir string, depth, limit uint) []*fileItem {
	if depth > limit {
		return nil
	}
	files, err := ioutil.ReadDir(filepath.Join(base, dir))
	if err != nil {
		panic(err)
	}

	var items []*fileItem
	for _, file := range files {
		rel := filepath.Join(dir, file.Name())
		items = append(items, &fileItem{base, rel, file})
		if file.IsDir() {
			items = append(items, dirwalk(base, rel, depth+1, limit)...)
			continue
		}
	}
	return items
}

func innerReadAllArticles() {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	aa := make([]*models.Article, 0)
	cc := make([]*models.Category, 0)
	ww := make(map[string][]string)
	var baseDir = beego.AppConfig.String("articles_dir")
	var items = dirwalk(baseDir, ".", 0, 1)
	var regArticleFile = regexp.MustCompile(`^(\d\d\d\d-\d\d-\d\d)_(\S+)\.md$`)
	var warn = func(item *fileItem, message string) {
		_, ok := ww[item.relative]
		if ok {
			ww[item.relative] = append(ww[item.relative], message)
		} else {
			ww[item.relative] = []string{message}
		}
		logs.Warn("[%s] %s", item.relative, message)
	}

	sort.Slice(items, func(i, j int) bool { return items[i].relative < items[j].relative })
	for _, item := range items {
		// parse slugs
		var filename string
		var categorySlug string
		splitted := strings.SplitN(item.relative, "/", 2)
		if len(splitted) == 1 {
			categorySlug = "-"
			filename = splitted[0]
		} else if len(splitted) == 2 {
			if splitted[0] == "-" {
				// skip "-/" dir
				warn(item, "Ignore `-` dir")
				continue
			}
			categorySlug = splitted[0]
			filename = splitted[1]
		} else {
			warn(item, "Ignore files with three levels of depth.")
			continue
		}

		filetype := FILE_TYPE_OTHER
		if regArticleFile.MatchString(filename) {
			filetype = FILE_TYPE_ARTICLE
		} else if filename == CATEGORY_FILE_NAME {
			filetype = FILE_TYPE_CATEGORY
		} else {
			continue
		}

		// start reading
		buf, err := ioutil.ReadFile(filepath.Join(baseDir, item.relative))
		if err != nil {
			warn(item, fmt.Sprintf("Failed to read: %s", err.Error()))
			continue
		}
		content := string(buf)

		switch filetype {
		case FILE_TYPE_ARTICLE:
			matched := regArticleFile.FindStringSubmatch(filename)
			if len(matched) != 3 {
				warn(item, "Invalid markdown file")
				continue
			}
			dateStr := matched[1]
			slug := matched[2]
			a := models.NewArticle()
			a.Title = slug
			a.Slug = slug
			a.CategorySlug = categorySlug
			a.Date = dateStr

			warning := a.LoadFromContent(content)
			if warning != "" {
				warn(item, fmt.Sprintf("Invalid header: %s", warning))
			}
			err = a.Validate()
			if err != nil {
				warn(item, fmt.Sprintf("Validation failed: %s", err.Error()))
				continue
			}
			a.Identify()
			aa = append(aa, a)
		case FILE_TYPE_CATEGORY:
			c := models.NewCategory()
			c.Slug = categorySlug
			c.Name = categorySlug
			err = c.FromJSON(content)
			if err != nil {
				warn(item, fmt.Sprintf("Failed to parse: %s", err.Error()))
			}
			c.Identify()
			cc = append(cc, c)
		default:
			// ignore
			continue
		}
	}
	SetCachedArticles(aa)
	SetCachedCategorys(cc)
	SetCachedWarnings(ww)
	logs.Info("Read %d As and %d Cs (%d warns)", len(aa), len(cc), len(ww))
}

func innerWriteArticle(a *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()

	if a.CategorySlug == "" {
		return fmt.Errorf("Category must not be empty: %+v", a)
	}
	if a.Slug == "" {
		return fmt.Errorf("Slug must not be empty: %+v", a)
	}

	articleDir := beego.AppConfig.String("articles_dir")
	baseDir := filepath.Join(articleDir, a.GetBaseDir())
	err := utils.EnsureDir(baseDir)
	if err != nil {
		return fmt.Errorf("Failed to mkdir: %s", err.Error())
	}

	mdPath := filepath.Join(articleDir, a.GetPath())
	if utils.FileExists(mdPath) { // file exists
		return fmt.Errorf("File `%s` does already exit.", a.GetPath())
	}

	content, err := a.ToText()
	if err != nil {
		return fmt.Errorf("Failed to serialize article: %s", err.Error())
	}
	err = ioutil.WriteFile(mdPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write article(%s): %s", mdPath, err.Error())
	}

	logs.Info("Succeeded to write article(`%s`)", mdPath)
	return nil
}

func innerRemoveArticle(a *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	if !a.Identified() {
		return fmt.Errorf("Removing article is not identified.")
	}

	articlesDir := beego.AppConfig.String("articles_dir")
	path := filepath.Join(articlesDir, a.GetPath())

	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("Failed to remove article file(%s): %s", path, err.Error())
	}

	logs.Info("Succeeded to remove article(`%s`)", path)
	return nil
}

func innerUpdateArticle(oldA, newA *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	if !oldA.Identified() {
		return fmt.Errorf("Old article is not identified.")
	}
	if newA.Identified() {
		return fmt.Errorf("New article is already identified.")
	}

	articlesDir := beego.AppConfig.String("articles_dir")

	oldPath := filepath.Join(articlesDir, oldA.GetPath())
	newPath := filepath.Join(articlesDir, newA.GetPath())

	newBaseDir := filepath.Join(articlesDir, newA.GetBaseDir())
	err := utils.EnsureDir(newBaseDir)
	if err != nil {
		return fmt.Errorf("Failed to create category dir: %s", err.Error())
	}

	if !utils.FileExists(oldPath) {
		return fmt.Errorf("Article `%s` does not exist in `%s`", oldA.JointedSlug(), newPath)
	}

	// if needed to move file
	fileChanged := oldPath != newPath
	if fileChanged {
		if utils.FileExists(newPath) {
			return fmt.Errorf("Already file exists in: %s", newPath)
		}
	}

	err = utils.EnsureDir(newBaseDir)
	if err != nil {
		return fmt.Errorf("Failed to ensure dir: %s", err.Error())
	}

	// 1. move file if needed
	if fileChanged {
		err = os.Rename(oldPath, newPath)
		if err != nil {
			return fmt.Errorf(
				"Failed to move file from `%s` to `%s`: %s",
				oldPath, newPath, err.Error())
		}
	}

	content, err := newA.ToText()
	if err != nil {
		return fmt.Errorf("Failed to serialize article: %s", err.Error())
	}

	err = ioutil.WriteFile(newPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write article(%s): %s", newPath, err.Error())
	}

	logs.Info("Succeeded to update article(`%s` -> `%s`)", oldPath, newPath)
	return nil
}

func WriteArticle(a *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch <- innerWriteArticle(a)
		ioWaiter.Done()
	}()
}

func ReadAllArticles() {
	ioWaiter.Add(1)
	go func() {
		innerReadAllArticles()
		ioWaiter.Done()
	}()
}

func RemoveArticle(a *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch <- innerRemoveArticle(a)
		ioWaiter.Done()
	}()
}

func UpdateArticle(oldA, newA *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch <- innerUpdateArticle(oldA, newA)
		ioWaiter.Done()
	}()
}
