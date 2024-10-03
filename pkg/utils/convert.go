package utils

import (
	"encoding/hex"
	"encoding/json"
)

func ConvertMapToString(m map[string]string) (string, error) {

	verifyInfoJsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	// convert bytes to string
	return hex.EncodeToString(verifyInfoJsonBytes), nil

}

func ConvertStringToMap(s string) (map[string]string, error) {

	verifyInfoBytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]string)

	err = json.Unmarshal(verifyInfoBytes, &mp)
	if err != nil {
		return nil, err
	}

	return mp, nil

}

func ConvertMapStringStringToMapStringAny(m map[string]string) (n map[string]any) {

	n = make(map[string]any)
	for k := range m {
		n[k] = m[k]
	}

	return
}
