package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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

func parseJson(input io.ReadCloser) map[string]interface{} {
	var body map[string]interface{}
	defer input.Close()
	if data, err := ioutil.ReadAll(input); err != nil {
		log.Fatalf("expected error to be nil got %v", err)
	} else if err := json.Unmarshal(data, &body); err != nil {
		log.Fatalf("Not a JSON, got %v", err)
	}
	return body
}

func parseJsonEsFile(input io.ReadCloser) esFile {
	var body esFile
	defer input.Close()
	if data, err := ioutil.ReadAll(input); err != nil {
		log.Fatalf("expected error to be nil got %v", err)
	} else if err := json.Unmarshal(data, &body); err != nil {
		log.Fatalf("Not a JSON, got %v", err)
	}
	return body
}
