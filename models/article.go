package models

import (
	// "encoding/json"
	// "errors"
	"fmt"
	"time"
	"strings"
	// "regexp"
	// "strconv"
	"github.com/goccy/go-yaml"

	"github.com/thedevsaddam/govalidator"
)


const HEADER_DELIMITTER = "---"

type ArticleHeader struct {
	Title string     `json:"title" yaml:",omitempty"`
	Tags []string    `json:"tags" yaml:",omitempty"`
	Aliases []string `json:"aliases" yaml:",omitempty"`
	Digest string    `json:"digest" yaml:",omitempty"`
	Image string     `json:"image" yaml:",omitempty"`
	Private bool     `json:"private" yaml:",omitempty"`
	Special bool     `json:"special" yaml:",omitempty"`
	Priority int     `json:"priority" yaml:",omitempty"`
}

type Article struct {
	Base
	ArticleHeader
	CategorySlug string `json:"categorySlug"`
	Slug string         `json:"slug"`
	Date string         `json:"date"`
	Body string         `json:"body"`
}

func init() {
	govalidator.AddCustomRule("strict_date_str", func(field, rule, message string, value interface{}) error {
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
	a.CategorySlug = "-"
	return &a
}

func (a *Article) GetBaseDir() string {
	if a.CategorySlug == "-" {
		return ""
	}
	return a.CategorySlug
}

func (a *Article) GetPath() string {
	base := fmt.Sprintf("%s_%s.md", a.Date, a.Slug)
	if a.CategorySlug == "-" {
		return base
	}
	return a.CategorySlug + "/" + base
}

func (a *Article) JointedSlug() string {
	return a.CategorySlug + "/" + a.Slug
}


func (a *Article) LoadFromContent(content string) string {
	// var header []string
	lines := strings.Split(content, "\n")
	hasHeaderStart := lines[0] == HEADER_DELIMITTER
	headerEndingLine := -1
	// header may exist
	if (hasHeaderStart) {
		for i, line := range lines[1:] {
			if line == HEADER_DELIMITTER {
				headerEndingLine = i + 1
			}
		}
	}
	// not (starting and ending)
	if !(hasHeaderStart && headerEndingLine > 0) {
		return ""
	}
	body := strings.Join(lines[headerEndingLine+1:len(lines)], "\n")
	headerText := strings.Join(lines[1:headerEndingLine], "\n")
	err := yaml.Unmarshal([]byte(headerText), &a.ArticleHeader)
	if err != nil {
		a.Body = content
		return fmt.Sprintf("Failed to parse header: %s", err.Error())
	}
	a.Body = body
	return ""
}

func (a *Article) Validate() error {
	rules := govalidator.MapData{
		"slug": []string{"required"},
		"date": []string{"required", "strict_date_str"},
		"categorySlug": []string{"required"},
	}

	opts := govalidator.Options{
		Data:  a,
		Rules: rules,
	}

	v := govalidator.New(opts)

	messages := v.ValidateStruct()
	if len(messages) > 0 {
		err := &ValidationError{Messages: messages}
		return err
	}
	return nil
}

func (a *Article) ToText() (string, error) {
	buf, err := yaml.Marshal(&a.ArticleHeader)
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
