package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

const (
	defaultType     = "local"
	defaultTemplate = "{{quote}}"
)

// Metadata describes the configuration of a story
type Metadata struct {
	Type     string
	Template string
}

// Story is metadata plus a set of variable chunks
type Story struct {
	Meta        Metadata
	Data        map[string]interface{}
	TypeObj     storyType `yaml:"-"`
	templateObj *template.Template
}

// NewStoryFromPath loads a new story generator from a file or URL
func NewStoryFromPath(path string) (Story, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return NewStoryFromURL(path)
	}
	return NewStoryFromFile(path)
}

// NewStoryFromFile loads a new Story generator from a config file
func NewStoryFromFile(filePath string) (Story, error) {
	s := Story{}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(file, &s)
	return s, err
}

// NewStoryFromURL loads a new Story generator from a URL to a config file
func NewStoryFromURL(url string) (Story, error) {
	s := Story{}

	resp, err := http.Get(url)
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(body, &s)
	return s, err
}

// Generate creates a story string
func (s *Story) Generate() (string, error) {
	if s.Meta.Type == "" {
		s.Meta.Type = defaultType
	}
	if s.Meta.Template == "" {
		s.Meta.Template = defaultTemplate
	}
	if s.TypeObj == nil {
		storyFunc, ok := storyTypes[s.Meta.Type]
		if !ok {
			return "", fmt.Errorf("Type not supported: %s", s.Meta.Type)
		}
		s.TypeObj = storyFunc()
	}
	if s.templateObj == nil {
		funcMap, err := s.TypeObj.Funcs(s)
		if err != nil {
			return "", err
		}
		s.templateObj, err = template.New("").Funcs(funcMap).Parse(s.Meta.Template)
		if err != nil {
			return "", err
		}
	}

	var result bytes.Buffer
	err := s.templateObj.Execute(&result, "")
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

type storyType interface {
	Funcs(*Story) (template.FuncMap, error)
}

var storyTypes = map[string]func() storyType{
	"local": func() storyType { return localStory{} },
	"s3":    func() storyType { return s3Story{} },
}
