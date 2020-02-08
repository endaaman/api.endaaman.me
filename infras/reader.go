package infras

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego"
)


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
        paths = append(paths, filepath.Join(dir, file.Name()))
    }

    return paths
}


func ReadAllArticles(ch chan<- bool) {
	ch <- true
}
