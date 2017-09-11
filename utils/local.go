package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"text/template"
)

type localStory struct {
}

func (ls localStory) Funcs(s *Story) (template.FuncMap, error) {
	funcMap := template.FuncMap{}
	for name, iface := range s.Data {
		list, err := ifaceToStringSlice(iface)
		if err != nil {
			return funcMap, err
		}
		funcMap[name] = localRandomizer(list)
	}
	return funcMap, nil
}

func ifaceToStringSlice(i interface{}) ([]string, error) {
	var list []string

	iSlice, ok := i.([]interface{})
	if !ok {
		return list, fmt.Errorf("YAML section is not a list")
	}

	for _, iStr := range iSlice {
		str, ok := iStr.(string)
		if !ok {
			return list, fmt.Errorf("YAML entry not a string")
		}
		list = append(list, str)
	}
	return list, nil
}

func getRandomIndex(max int) (int, error) {
	maxBig := int64(max)
	n, err := rand.Int(rand.Reader, big.NewInt(maxBig))
	if err != nil {
		return 0, err
	}
	i := int(n.Int64())
	return i - 1, nil
}

func localRandomizer(list []string) func() (string, error) {
	return func() (string, error) {
		index, err := getRandomIndex(len(list))
		if err != nil {
			return "", err
		}
		return list[index], nil
	}
}
