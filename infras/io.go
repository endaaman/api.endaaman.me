package infras

import (
    "os"
    "regexp"
    "fmt"
    "strings"
    "sync"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/models"
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

func dirwalk(dir string, depth, limit uint) []string {
	if depth > limit {
		return nil
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			paths = append(paths, dirwalk(path, depth + 1, limit)...)
			continue
		}
		paths = append(paths, path)
    }
    return paths
}

func innerReadAllArticles() {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	aa := make([]*models.Article, 0)
	cc := make([]*models.Category, 0)
	ww := make(map[string][]string)
	var baseDir = beego.AppConfig.String("articles_dir")
	var paths = dirwalk(baseDir, 0, 1)
	var regMd = regexp.MustCompile(`\.md$`)
	var warn = func(path, message string) {
		_, ok := ww[path]
		if ok {
			ww[path] = append(ww[path], message)
		} else {
			ww[path] = []string{message}
		}
		logs.Warn("[%s] %s", path, message)
	}
	for _, path := range paths {
		// compute rel
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			warn("common", fmt.Sprintf("Failed compute rel: %s", err.Error()))
			continue
		}

		// parse slugs
		var filename string
		var categorySlug string
		splitted := strings.SplitN(rel, "/", 2)
		if len(splitted) == 1 {
			categorySlug = "-"
			filename = splitted[0]
		} else if len(splitted) == 2 {
			if splitted[0] == "-" {
				// skip "-/" dir
				continue
			}
			categorySlug = splitted[0]
			filename = splitted[1]
		} else {
			warn(path, "Invalid path")
			continue
		}

		// ignore unnecessary
		ignored := true
		filetype := FILE_TYPE_OTHER
		if regMd.MatchString(filename) {
			filetype = FILE_TYPE_ARTICLE
			ignored = false
		}
		if filename == CATEGORY_FILE_NAME {
			filetype = FILE_TYPE_CATEGORY
			ignored = false
		}
		if ignored {
			continue
		}

		// start reading
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			warn(path, fmt.Sprintf("Failed to read: %s", err.Error()))
			continue
		}
		content := string(buf)

		switch filetype {
		case FILE_TYPE_ARTICLE:
			stat, err := os.Stat(path)
			if err != nil {
				warn(path, fmt.Sprintf("Failed to stat: %s", err.Error()))
				continue
			}
			slug := regMd.ReplaceAllString(filename, "")
			dateStr := stat.ModTime().Format("2006-01-02")
			a := models.NewArticle()
			a.Title = slug
			a.Slug = slug
			a.CategorySlug = categorySlug
			a.Date = dateStr
			header, body, err := models.SplitArticleHeaderAndBody(content)
			if err != nil {
				warn(path, fmt.Sprintf("Failed to parse markdown: %s", err.Error()))
			}
			if header != nil {
				a.ArticleHeader = *header
			}
			a.Body = body
			err = a.Validate()
			if err != nil {
				warn(path, fmt.Sprintf("Invalid header: %s", err.Error()))
			}
			a.Identify()
			aa = append(aa, a)
		case FILE_TYPE_CATEGORY:
			c := models.NewCategory()
			c.Slug = categorySlug
			c.Name = categorySlug
			err = c.FromJSON(content)
			if err != nil {
				warn(path, fmt.Sprintf("Failed to parse: %s", err.Error()))
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

func WriteArticle(a *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	ioMutex.Lock()
	defer ioMutex.Unlock()
	defer ioWaiter.Done()

	if a.CategorySlug == "" {
		ch<- fmt.Errorf("Category must not be empty: %+v", a)
		return
	}
	if a.Slug == "" {
		ch<- fmt.Errorf("Slug must not be empty: %+v", a)
		return
	}

	articleDir := beego.AppConfig.String("articles_dir")
	var categoryDir string
	if a.CategorySlug == "-" {
		categoryDir = ""
	} else {
		categoryDir = a.CategorySlug
	}
	baseDir := filepath.Join(articleDir, categoryDir)
	err := os.MkdirAll(baseDir, 0777);
    if err != nil {
		ch<- fmt.Errorf("Failed to mkdir: %s", err.Error())
		return
    }

	mdPath := filepath.Join(baseDir, a.Slug + ".md")
    _, err = os.Stat(mdPath)
	if err == nil { // file exists
		ch<- fmt.Errorf("Already `%s/%s` does already exit.", a.CategorySlug, a.Slug)
		return
	}

	content, err := a.ToText()
    if err != nil {
		ch<- fmt.Errorf("Failed to serialize article: %s", err.Error())
		return
    }
	err = ioutil.WriteFile(mdPath, []byte(content), 0644)
    if err != nil {
		ch<- fmt.Errorf("Failed to write article(%s): %s", mdPath, err.Error())
		return
    }
	logs.Info("Success wrote article(`%s`)", mdPath)
	ch<- nil
	return
}

func innerRemoveArticle(a *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	if (!a.Identified()) {
		return fmt.Errorf("Removing article is not identified.")
	}
	// TODO: impl delete
	return nil
}

func innerReplaceArticle(oldA, newA *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	if (!oldA.Identified()) {
		return fmt.Errorf("Old article is not identified.")
	}
	if (newA.Identified()) {
		return fmt.Errorf("New article is already identified.")
	}
	// TODO: impl delete and create
	return nil
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
		ch<- innerRemoveArticle(a)
		ioWaiter.Done()
	}()
}

func ReplaceArticle(oldA, newA *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch<- innerReplaceArticle(oldA, newA)
		ioWaiter.Done()
	}()
}
