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
var regMd = regexp.MustCompile(`\.md$`)


// func GetReader() *sync.WaitGroup {
// 	return ioWaiter
// }

func WaitIO() {
	ioWaiter.Wait()
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		name := file.Name()
		if regMd.MatchString(name) {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
    }
    return paths
}

func innerReadAllArticles() {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	ww := make([]string, 0)
	aa := make([]*models.Article, 0)
	var errCount = 0
	var warningCount = 0
	var baseDir = beego.AppConfig.String("articles_dir")
	var paths = dirwalk(baseDir)

	var save_err = func(path string, err error) {
		w := fmt.Sprintf("%s: %s", path, err.Error())
		ww = append(ww, w)
		logs.Warn("Article error: %s", w)
		errCount += 1
	}

	for _, path := range paths {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			save_err(path, err)
			continue
		}

		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			save_err(path, err)
			continue
		}
		slug := regMd.ReplaceAllString(rel, "")
		category := ""
		splitted := strings.SplitN(slug, "/", 2)
		if len(splitted) == 2 {
			category = splitted[0]
			slug = splitted[1]
			if category == "-" {
				// skip "-" for article may duplicate
				continue
			}
		}

		a := models.NewArticle()
		content := string(buf)
		headerLoaded := a.FromText(content, category, slug, fi.ModTime().Format("2006-01-02"))
		if !headerLoaded {
			w := fmt.Sprintf("%s: failed to parse header", path)
			ww = append(ww, w)
			logs.Warn("Article warning: %s", w)
			logs.Warn("Content: %s", content)
			warningCount += 1
		}
		a.Identify()
		aa = append(aa, a)
	}
	SetCachedArticles(aa)
	SetCachedWarnings(ww)
	logs.Info("Read %d articles (%d warns, %d errs).", len(paths), warningCount, errCount)
}

func innerWriteArticle(a *models.Article) error {
	ioMutex.Lock()
	defer ioMutex.Unlock()
	if a.Category == "" {
		return fmt.Errorf("Category must not be empty: %+v", a)
	}
	if a.Slug == "" {
		return fmt.Errorf("Slug must not be empty: %+v", a)
	}

	articleDir := beego.AppConfig.String("articles_dir")
	var categoryDir string
	if a.Category == "-" {
		categoryDir = ""
	} else {
		categoryDir = a.Category
	}
	baseDir := filepath.Join(articleDir, categoryDir)
	err := os.MkdirAll(baseDir, 0777);
    if err != nil {
		return fmt.Errorf("Failed to mkdir: %s", err.Error())
    }

	mdPath := filepath.Join(baseDir, a.Slug + ".md")
    _, err = os.Stat(mdPath)
	if err == nil { // file exists
		return fmt.Errorf("Already `%s/%s` does already exit.", a.Category, a.Slug)
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
