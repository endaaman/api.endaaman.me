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

var ioWaiter = new(sync.WaitGroup)
var ioMutex = new(sync.Mutex)


// func GetReader() *sync.WaitGroup {
// 	return ioWaiter
// }

const CATEGORY_FILE_NAME = "meta.json"
const (
	FILE_TYPE_ARTICLE = iota
	FILE_TYPE_CATEGORY
	FILE_TYPE_OTHER
)

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
	var baseDir = beego.AppConfig.String("articles_dir")
	var paths = dirwalk(baseDir, 0, 1)
	var regMd = regexp.MustCompile(`\.md$`)
	for _, path := range paths {
		// compute rel
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			logs.Error("Failed compute rel: %s", err.Error())
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
			logs.Error("Invalid path: %s", path)
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
			logs.Error("Failed to read `%s`: %s", path, err.Error())
			continue
		}
		content := string(buf)

		switch filetype {
		case FILE_TYPE_ARTICLE:
			stat, err := os.Stat(path)
			if err != nil {
				logs.Error("Failed to stat `%s`: %s", path, err.Error())
				continue
			}
			slug := regMd.ReplaceAllString(filename, "")
			dateStr := stat.ModTime().Format("2006-01-02")
			a := models.NewArticle()
			a.FromText(categorySlug, slug, dateStr, content)
			a.Identify()
			aa = append(aa, a)
		case FILE_TYPE_CATEGORY:
			c := models.NewCategory()
			c.Identify()
			cc = append(cc, c)
		default:
			// ignore
			continue
		}
	}
	SetCachedArticles(aa)
	SetCachedCategorys(cc)
	logs.Info("Read %d As and %d Cs", len(aa), len(cc))
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
	var categoryDir string
	if a.CategorySlug == "-" {
		categoryDir = ""
	} else {
		categoryDir = a.CategorySlug
	}
	baseDir := filepath.Join(articleDir, categoryDir)
	err := os.MkdirAll(baseDir, 0777);
    if err != nil {
		return fmt.Errorf("Failed to mkdir: %s", err.Error())
    }

	mdPath := filepath.Join(baseDir, a.Slug + ".md")
    _, err = os.Stat(mdPath)
	if err == nil { // file exists
		return fmt.Errorf("Already `%s/%s` does already exit.", a.CategorySlug, a.Slug)
	}

	content, err := a.ToText()
    if err != nil {
		return fmt.Errorf("Failed to serialize article: %s", err.Error())
    }
	err = ioutil.WriteFile(mdPath, []byte(content), 0644)
    if err != nil {
		return fmt.Errorf("Failed to write article(%s): %s", mdPath, err.Error())
    }
	logs.Info("Success wrote article(`%s`)", mdPath)
	return nil
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

func WriteArticle(a *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch<- innerWriteArticle(a)
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
