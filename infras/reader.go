package infras

import (
    "os"
    "regexp"
    "fmt"
    "sync"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
	"github.com/endaaman/api.endaaman.me/models"
)

var reader = &sync.WaitGroup{}
var reg_md = regexp.MustCompile(`\.md$`)


func GetReader() *sync.WaitGroup {
	return reader
}

func WaitReader() {
	reader.Wait()
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

func ReadAllArticles() *sync.WaitGroup {
	reader.Add(1)

	go func() {
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
			a := models.Article{}
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

			a.FromText(string(buf), slug, fi.ModTime().Format("2006-01-02"))
			if a.Warning != "" {
				w := fmt.Sprintf("%s: %s", path, a.Warning)
				ww = append(ww, w)
				logs.Warn("Article warning: %s", w)
				warningCount += 1
			}
			aa = append(aa, &a)
		}
		SetCachedArticles(aa)
		SetCachedWarnings(ww)
		reader.Done()
		logs.Info("Read %d items (%d warns, %d errs).", len(paths), warningCount, errCount)
	}()
	return reader
}
