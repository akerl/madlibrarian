package utils

import (
	"fmt"
	"text/template"
)

type s3Story struct {
}

func (ss s3Story) Funcs(s *Story) (template.FuncMap, error) {
	funcMap := template.FuncMap{}

	s3InfoIface, ok := s.Data["s3"]
	if !ok {
		return funcMap, fmt.Errorf("No S3 section in YAML")
	}
	s3Info, ok := s3InfoIface.(map[string]string)
	if !ok {
		return funcMap, fmt.Errorf("S3 section is not a map of strings to strings")
	}

	chunksIface, ok := s.Data["chunks"]
	if !ok {
		return funcMap, fmt.Errorf("Chunks not defined")
	}
	chunks, err := ifaceToStringSlice(chunksIface)
	if err != nil {
		return funcMap, err
	}

	for _, name := range chunks {
		funcMap[name] = s3Randomizer(s3Info, name)
	}
	return funcMap, nil
}

func s3Randomizer(s3Info map[string]string, chunk string) func() (string, error) {
	return func() (string, error) {
		// TODO: Implement randomizer for reading from S3
		return "PLACEHOLDER", nil
	}
}
