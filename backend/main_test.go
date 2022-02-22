package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func toJson(v interface{}) io.Reader {
	if bytes_, err := json.Marshal(v); err == nil {
		return bytes.NewReader(bytes_)
	} else {
		panic(err)
	}
}

func TestAddFileTxt(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Post(fmt.Sprintf("%s/api/file/upload", ts.URL), "application/json", toJson(gin.H{
		"name": "file.txt",
		"size": 100,
		"tags": []string{"text", "document"},
	}))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	val, ok := resp.Header["Content-Type"]

	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}
