package infras

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/config"
	"github.com/endaaman/api.endaaman.me/models"
	"github.com/endaaman/api.endaaman.me/utils"
)

const WARNING_KEY_ROOT = "_ROOT"

const META_FILE_NAME = "meta.json"
const (
	FILE_TYPE_ARTICLE = iota
	FILE_TYPE_CATEGORY
	FILE_TYPE_OTHER
)

type CategoryMetaEntry struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type FileItem struct {
	rel  string
	file os.FileInfo
}

var ioWaiter = new(sync.WaitGroup)
var ioMutex = new(sync.Mutex)

func WaitIO() {
	ioWaiter.Wait()
}

func dirwalk(base, dir string, depth, limit uint) []*FileItem {
	if depth > limit {
		return nil
	}
	files, err := ioutil.ReadDir(filepath.Join(base, dir))
	if err != nil {
		panic(err)
	}

	var items []*FileItem
	for _, file := range files {
		rel := filepath.Join(dir, file.Name())
		if file.IsDir() {
			items = append(items, dirwalk(base, rel, depth+1, limit)...)
		}
		item := FileItem{
			rel:  rel,
			file: file,
		}
		items = append(items, &item)
	}
	return items
}

func appendWarning(ww map[string][]string, item string, message string) {
	_, ok := ww[item]
	if ok {
		ww[item] = append(ww[item], message)
	} else {
		ww[item] = []string{message}
	}
	logs.Warn("[%s] %s", item, message)
}

func loadArticles(items []*FileItem, ww map[string][]string) []*models.Article {
	articlesDir := config.GetArticlesDir()
	regMarkdown := regexp.MustCompile(`^\S+\.md$`)
	regArticleFile := regexp.MustCompile(`^(\d\d\d\d-\d\d-\d\d)_(\S+)\.md$`)

	aa := make([]*models.Article, 0)
	for _, item := range items {
		var filename string
		var categorySlug string
		splitted := strings.SplitN(item.rel, "/", 2)
		if len(splitted) == 1 {
			categorySlug = "-"
			filename = splitted[0]
		} else if len(splitted) == 2 {
			if splitted[0] == "-" {
				// skip "-/" dir
				logs.Debug("Ignore `-` dir: %s", item)
				continue
			}
			categorySlug = splitted[0]
			filename = splitted[1]
		} else {
			// this should be never reached
			continue
		}

		matched := regArticleFile.FindStringSubmatch(filename)
		if len(matched) != 3 {
			if regMarkdown.MatchString(filename) {
				appendWarning(ww, item.rel, "Invalid markdown file")
			}
			continue
		}

		buf, err := ioutil.ReadFile(filepath.Join(articlesDir, item.rel))
		if err != nil {
			appendWarning(ww, item.rel, fmt.Sprintf("Failed to read file: %s", err.Error()))
			continue
		}
		content := string(buf)

		dateStr := matched[1]
		slug := matched[2]
		a := models.NewArticle()
		a.Title = slug
		a.Slug = slug
		a.CategorySlug = categorySlug
		a.Date = dateStr

		warning := a.LoadFromContent(content)
		if warning != "" {
			appendWarning(ww, item.rel, fmt.Sprintf("Invalid header: %s", warning))
		}
		err = a.Validate()
		if err != nil {
			appendWarning(ww, item.rel, fmt.Sprintf("Validation failed: %s", err.Error()))
			continue
		}
		a.Identify()
		aa = append(aa, a)
	}
	return aa
}

func loadCategories(items []*FileItem, ww map[string][]string) []*models.Category {
	articlesDir := config.GetArticlesDir()
	metaPath := filepath.Join(articlesDir, META_FILE_NAME)
	cc := make([]*models.Category, 0)
	meta := make(map[string]CategoryMetaEntry, 0)

	if utils.FileExists(metaPath) {
		buf, err := ioutil.ReadFile(metaPath)
		if err != nil {
			appendWarning(ww, META_FILE_NAME, fmt.Sprintf("Failed to read: %s", err.Error()))
		}
		err = json.Unmarshal(buf, &meta)
		if err != nil {
			appendWarning(ww, META_FILE_NAME, fmt.Sprintf("Invalid json: %s", err.Error()))
		}
	}

	slugs := make([]string, 1)
	slugs[0] = "-" // default category
	for _, item := range items {
		if !item.file.IsDir() {
			continue
		}
		slugs = append(slugs, item.rel)
	}

	for _, slug := range slugs {
		c := models.NewCategory(slug)
		entry, ok := meta[slug]
		if ok {
			c.Name = entry.Name
			c.Priority = entry.Priority
		}
		c.Identify()
		cc = append(cc, c)
	}

	return cc
}

func innerReadAllArticles() {
	ioMutex.Lock()
	defer ioMutex.Unlock()

	aa := make([]*models.Article, 0)
	cc := make([]*models.Category, 0)
	ww := make(map[string][]string)

	articlesDir := config.GetArticlesDir()
	if utils.IsDir(articlesDir) {
		items := dirwalk(articlesDir, ".", 0, 1)
		sort.Slice(items, func(i, j int) bool { return items[i].rel < items[j].rel })

		aa = loadArticles(items, ww)
		cc = loadCategories(items, ww)

		sort.Slice(aa, func(a, b int) bool { return aa[a].Compare(aa[b]) })
		sort.Slice(cc, func(a, b int) bool { return cc[a].Compare(cc[b]) })
	} else {
		appendWarning(ww, WARNING_KEY_ROOT, fmt.Sprintf("Articles Dir(%s) does not exists", articlesDir))
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

	articlesDir := config.GetArticlesDir()
	baseDir := filepath.Join(articlesDir, a.GetBaseDir())
	err := utils.EnsureDir(baseDir)
	if err != nil {
		return fmt.Errorf("Failed to mkdir: %s", err.Error())
	}

	mdPath := filepath.Join(articlesDir, a.GetPath())
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

	articlesDir := config.GetArticlesDir()
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

	articlesDir := config.GetArticlesDir()

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
