package main

import (
	"bytes"
	"encoding/json"
	"io"
)

func toJson(v interface{}) io.Reader {
	if bytesData, err := json.Marshal(v); err == nil {
		return bytes.NewReader(bytesData)
	} else {
		panic(err)
	}
}

func stringKeys(v map[string]bool) []string {
	res := make([]string, 0)
	for k := range v {
		res = append(res, k)
	}
	return res
}
