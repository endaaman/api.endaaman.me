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
var reg_md = regexp.MustCompile(`\.md$`)


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
		if reg_md.MatchString(name) {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
    }
    return paths
}

func innerReadAllArticles() {
	var ww []string
	var aa []*models.Article
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
		slug := reg_md.ReplaceAllString(rel, "")
		category := ""
		splitted := strings.SplitN(slug, "/", 2)
		if len(splitted) == 2 {
			slug = splitted[1]
			category = splitted[1]
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
	logs.Info("Read %d items (%d warns, %d errs).", len(paths), warningCount, errCount)
}

func ReadAllArticles() *sync.WaitGroup {
	ioWaiter.Add(1)
	go func() {
		innerReadAllArticles()
		ioWaiter.Done()
	}()
	return ioWaiter
}

func innerWriteArticle(a *models.Article) {
	baseDir := beego.AppConfig.String("articles_dir")
	err := os.MkdirAll(filepath.Join(baseDir, a.Category), 0777);
    if err != nil {
		fmt.Println("mkdir:", err)
		return
    }

	mdPath := filepath.Join(baseDir, a.Category, a.Slug + ".md")
	content := a.ToText()
	err = ioutil.WriteFile(mdPath, []byte(content), 0644)
    if err != nil {
		fmt.Println("write:", err)
		return
    }
	fmt.Println("write done")
}

func WriteArticle(a *models.Article) *sync.WaitGroup {
	ioWaiter.Add(1)
	go func() {
		innerWriteArticle(a)
		ioWaiter.Done()
	}()
	return ioWaiter
}

