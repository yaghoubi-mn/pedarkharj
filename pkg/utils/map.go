package utils

import (
	"fmt"
)

func CheckMapHaveKeys[C comparable, A any](m map[C]A, keys ...C) error {

	// find keys
	keys2 := make([]C, 0, len(m))
	for k := range m {
		keys2 = append(keys2, k)
	}

	var isExist bool
	for _, k := range keys {
		isExist = false
		for _, k2 := range keys2 {
			if k2 == k {
				isExist = true
				break
			}

		}

		if !isExist {
			return fmt.Errorf("key \"%v\" not found", k)
		}
	}

	return nil

}
