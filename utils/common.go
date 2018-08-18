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
	iSlice, ok := i.([]interface{})
	if !ok {
		return []string{}, fmt.Errorf("YAML section is not a list")
	}

	list := make([]string, len(iSlice))
	for index, iStr := range iSlice {
		str, ok := iStr.(string)
		if !ok {
			return list, fmt.Errorf("YAML entry not a string")
		}
		list[index] = str
	}
	return list, nil
}
