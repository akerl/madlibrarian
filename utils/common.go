package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func getRand(max int) (int, error) {
	maxBig := int64(max)
	n, err := rand.Int(rand.Reader, big.NewInt(maxBig))
	if err != nil {
		return 0, err
	}
	i := int(n.Int64())
	return i, nil
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
