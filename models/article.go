package models

import (
	"fmt"
	// "errors"
	// "encoding/json"
	"time"
	"strings"
	// "regexp"
	// "strconv"
	"github.com/goccy/go-yaml"

	"github.com/thedevsaddam/govalidator"
)


const HEADER_DELIMITTER = "---"

type Header struct {
	Title string     `json:"title" yaml:",omitempty"`
	Tags []string    `json:"tags" yaml:",omitempty"`
	Aliases []string `json:"aliases" yaml:",omitempty"`
	Digest string    `json:"digest" yaml:",omitempty"`
	Image string     `json:"image" yaml:",omitempty"`
	Private bool     `json:"private" yaml:",omitempty"`
	Special bool     `json:"special" yaml:",omitempty"`
	Priority int     `json:"priority" yaml:",omitempty"`
	Date string      `json:"date" yaml:",omitempty"`
}


type Article struct {
	Header
	Slug string     `json:"slug"`
	Category string `json:"category"`
	Body string     `json:"body"`
	Warning string  `json:"warning"`
}

func init() {
	govalidator.AddCustomRule("strict_date_str", func(field string, rule string, message string, value interface{}) error {
		s := value.(string)
		if s == "" {
			return nil
		}
		_, err := time.Parse("2006-01-02", s)
		if err != nil {
			return fmt.Errorf("`%s` is not valid date", s)
		}
		return nil
	})
}

func NewArticle() *Article {
	a := Article{}
	a.Tags = make([]string, 0)
	a.Aliases = make([]string, 0)
	a.Date = time.Now().Format("2006-01-02")
	a.Category = "-"
	return &a
}

func (a *Article) FromText(text string, category string, slug string, date string) {
	a.Title = slug
	a.Slug = slug
	if category == "" {
		a.Category = "-"
	} else {
		a.Category = category
	}
	a.Date = date

	// var header []string
	lines := strings.Split(text, "\n")
	hasHeaderStart := lines[0] == HEADER_DELIMITTER
	headerEndingLine := -1
	if (hasHeaderStart) {
		for i, line := range lines[1:] {
			if line == HEADER_DELIMITTER {
				headerEndingLine = i + 1
			}
		}
	}
	if !(hasHeaderStart && headerEndingLine > 0) {
		a.Body = text
		return
	}

	a.Body = strings.Join(lines[headerEndingLine+1:len(lines)], "\n")
	headerText := strings.Join(lines[1:headerEndingLine], "\n")
	err := yaml.Unmarshal([]byte(headerText), &a.Header)
	if err != nil {
		a.Warning = err.Error()
	}
}

func (a *Article) Validate() map[string][]string {
	rules := govalidator.MapData{
		"slug": []string{"required"},
		"date": []string{"required", "strict_date_str"},
		"category": []string{"required"},
	}

	opts := govalidator.Options{
		Data:  a,
		Rules: rules,
	}

	v := govalidator.New(opts)

	e := v.ValidateStruct()
	if len(e) > 0 {
		return e
	}
	return nil
}

func (a *Article) ToText() (string, error) {
	buf, err := yaml.Marshal(&a.Header)
	if err != nil {
		return "", err
	}
	headerText := string(buf)

	content := a.Body
	if headerText != "{}\n" {
		template := "---\n%s---\n%s"
		content = fmt.Sprintf(template, headerText, a.Body)
	}
	return content, nil
}
