package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

const (
	defaultType     = "local"
	defaultTemplate = "{{quote}}"
)

var urlPrefixes = []string{"http://", "https://"}

// Metadata describes the configuration of a story
type Metadata struct {
	Type     string
	Template string
}

// Story is metadata plus a set of variable chunks
type Story struct {
	Meta     Metadata
	Data     map[string]interface{}
	author   author
	template *template.Template
}

// NewStoryFromPath loads a new story generator from a file or URL
func NewStoryFromPath(path string) (Story, error) {
	for _, prefix := range urlPrefixes {
		if strings.HasPrefix(path, prefix) {
			return NewStoryFromURL(path)
		}
	}
	return NewStoryFromFile(path)
}

// NewStoryFromFile loads a new Story generator from a config file
func NewStoryFromFile(filePath string) (Story, error) {
	file, err := ioutil.ReadFile(filePath) // #nosec
	if err != nil {
		return Story{}, err
	}

	return NewStoryFromText(file)
}

// NewStoryFromURL loads a new Story generator from a URL to a config file
func NewStoryFromURL(url string) (Story, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Story{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Story{}, err
	}

	return NewStoryFromText(body)
}

// NewStoryFromText loads a new Story generator from a string
func NewStoryFromText(text []byte) (Story, error) {
	s := Story{}
	err := yaml.Unmarshal(text, &s)
	if err != nil {
		return s, err
	}
	err = s.Init()
	return s, err
}

// Init sets up initial state for the Story object
func (s *Story) Init() error {
	if s.Meta.Type == "" {
		s.Meta.Type = defaultType
	}
	if s.Meta.Template == "" {
		s.Meta.Template = defaultTemplate
	}
	if s.author == nil {
		authorFunc, ok := authorTypes[s.Meta.Type]
		if !ok {
			return fmt.Errorf("Type not supported: %s", s.Meta.Type)
		}
		s.author = authorFunc()
	}
	if s.template == nil {
		funcMap, err := s.author.Funcs(s)
		if err != nil {
			return err
		}
		s.template, err = template.New("").Funcs(funcMap).Parse(s.Meta.Template)
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate creates a story string
func (s *Story) Generate() (string, error) {
	var result bytes.Buffer
	err := s.template.Execute(&result, "")
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// Upload copies a story into S3
func (s *Story) Upload(bucket, prefix string) (*Story, error) {
	return s.author.Upload(s, bucket, prefix)
}

type author interface {
	Funcs(*Story) (template.FuncMap, error)
	Upload(*Story, string, string) (*Story, error)
}

var authorTypes = map[string]func() author{
	"local": func() author { return localAuthor{} },
	"s3":    func() author { return s3Author{} },
}
