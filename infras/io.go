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

var ioWaiter = &sync.WaitGroup{}
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
		a.FromText(string(buf), category, slug, fi.ModTime().Format("2006-01-02"))
		if a.Warning != "" {
			w := fmt.Sprintf("%s: %s", path, a.Warning)
			ww = append(ww, w)
			logs.Warn("Article warning: %s", w)
			warningCount += 1
		}
		aa = append(aa, a)
	}
	SetCachedArticles(aa)
	SetCachedWarnings(ww)
	logs.Info("Read %d articles (%d warns, %d errs).", len(paths), warningCount, errCount)
}

func ReadAllArticles() {
	ioWaiter.Add(1)
	go func() {
		innerReadAllArticles()
		ioWaiter.Done()
	}()
}

func innerWriteArticle(a *models.Article) error {
	if a.Category == "" {
		return fmt.Errorf("Category must not be empty: %+v", a)
	}
	if a.Slug == "" {
		return fmt.Errorf("Slug must not be empty: %+v", a)
	}

	articleDir := beego.AppConfig.String("articles_dir")
	var categoryDir string
	fmt.Println("CATEGORY: ", a.Category)
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

func WriteArticle(a *models.Article, ch chan<- error) {
	ioWaiter.Add(1)
	go func() {
		ch<- innerWriteArticle(a)
		ioWaiter.Done()
	}()
}

