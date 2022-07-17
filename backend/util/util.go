package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

type M map[string]interface{}

func ToJson(v interface{}) io.Reader {
	if bytesData, err := json.Marshal(v); err == nil {
		return bytes.NewReader(bytesData)
	} else {
		panic(err)
	}
}

/*func StringKeys(v map[string]bool) []string {
	res := make([]string, 0)
	for k := range v {
		res = append(res, k)
	}
	return res
}*/

func ParseJsonTyped(input io.ReadCloser, v interface{}) {
	defer input.Close()
	if data, err := ioutil.ReadAll(input); err != nil {
		log.Fatalf("expected error to be nil got %v", err)
	} else if err := json.Unmarshal(data, v); err != nil {
		log.Fatalf("Not a JSON, got %v", err)
	}
}

// TODO is there an official generified lib?
func Map[T any, R any](slice []T, f func(T) R) []R {
	res := make([]R, 0, len(slice))
	for _, t := range slice {
		res = append(res, f(t))
	}
	return res
}
