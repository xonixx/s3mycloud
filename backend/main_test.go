package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type M map[string]interface{}

func toJson(v interface{}) io.Reader {
	if bytesData, err := json.Marshal(v); err == nil {
		return bytes.NewReader(bytesData)
	} else {
		panic(err)
	}
}

func checkResponse(t *testing.T, resp *http.Response, err error, expectedStatus int) {
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != expectedStatus {
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

func readJsonAsMap(t *testing.T, resp *http.Response) M {
	var body M
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("expected error to be nil got %v", err)
	} else if err := json.Unmarshal(data, &body); err != nil {
		t.Errorf("Not a JSON, got %v", err)
	}
	return body
}

func checkAndReadRespJson(t *testing.T, resp *http.Response, err error, expectedStatus int) M {
	checkResponse(t, resp, err, expectedStatus)
	return readJsonAsMap(t, resp)
}

func getJsonField(t *testing.T, m M, field string) interface{} {
	val, ok := m[field]
	if !ok {
		t.Fatalf("should contain %s", field)
		return nil
	} else {
		return val
	}
}

func query(v interface{}, path ...string) interface{} {
	for _, f := range path {
		if m, ok := v.(map[string]interface{}); ok {
			v = m[f]
		} else if m, ok := v.(M); ok {
			v = m[f]
		} else if l, ok := v.([]interface{}); ok {
			if intVal, err := strconv.Atoi(f); err == nil {
				v = l[intVal]
			} else {
				panic("Not an int: " + f)
			}
		} else {
			panic(fmt.Sprintf("Not a map/list: %v", v))
		}
	}
	return v
}

type testHelper struct {
	t  *testing.T
	ts *httptest.Server
}

func (th testHelper) get(path string) M {
	resp, err := http.Get(fmt.Sprintf("%s/%s", th.ts.URL, path))
	return checkAndReadRespJson(th.t, resp, err, http.StatusOK)
}

func (th testHelper) assertEqualsJsonPath(json M, expectedVal interface{}, path ...string) {
	val := query(json, path...)
	if val != expectedVal {
		th.t.Fatalf("expected %v != actual %v", expectedVal, val)
	}
}

func (th testHelper) setupFiles() {
	cleanupMemStorage()

	for _, f := range []M{
		{
			"name": "aaa.mp3",
			"size": 1000_000,
			"tags": []string{"music", "pop"},
		},
		{
			"name": "bbb.txt",
			"size": 100,
			"tags": []string{"text", "document"},
		},
		{
			"name": "ccc",
			"size": 0,
		},
		{
			"name": "ddd",
			"size": 200,
		},
		{
			"name": "eee",
			"size": 500,
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/api/file/upload", th.ts.URL), "application/json", toJson(f))

		//fmt.Println("resp JSON:", readJsonAsMap(t, resp))

		respJson := checkAndReadRespJson(th.t, resp, err, http.StatusCreated)
		getJsonField(th.t, respJson, "id")
		getJsonField(th.t, respJson, "uploadUrl")
	}
}

func TestAddSingleFile(t *testing.T) {
	cleanupMemStorage()
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Post(fmt.Sprintf("%s/api/file/upload", ts.URL), "application/json", toJson(M{
		"name": "file.txt",
		"size": 100,
		"tags": []string{"text", "document"},
	}))

	//fmt.Println("resp JSON:", readJsonAsMap(t, resp))

	respJson := checkAndReadRespJson(t, resp, err, http.StatusCreated)

	id := getJsonField(t, respJson, "id").(string)
	//fmt.Println(id)
	getJsonField(t, respJson, "uploadUrl")

	resp, err = http.Get(fmt.Sprintf("%s/api/file", ts.URL))
	respJson = checkAndReadRespJson(t, resp, err, http.StatusOK)

	if total := int(getJsonField(t, respJson, "total").(float64)); total != 1 {
		t.Fatalf("wrong total: %d", total)
	}
	resultId := query(respJson, "page", "0", "id").(string)
	if id != resultId {
		t.Fatalf("resultId %s != id %s", resultId, id)
	}
}

func TestListDefaultSorting(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()
	th := testHelper{t, ts}
	th.setupFiles()

	respJson := th.get("api/file")

	if total := int(getJsonField(t, respJson, "total").(float64)); total != 5 {
		t.Fatalf("wrong total: %d", total)
	}

	th.assertEqualsJsonPath(respJson, "eee", "page", "0", "name")
	th.assertEqualsJsonPath(respJson, "ddd", "page", "1", "name")
	th.assertEqualsJsonPath(respJson, "ccc", "page", "2", "name")
	th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "3", "name")
	th.assertEqualsJsonPath(respJson, "aaa.mp3", "page", "4", "name")
}

func TestListOlderFirst(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()
	th := testHelper{t, ts}
	th.setupFiles()

	respJson := th.get("api/file?sort=uploaded,asc")

	if total := int(getJsonField(t, respJson, "total").(float64)); total != 5 {
		t.Fatalf("wrong total: %d", total)
	}

	th.assertEqualsJsonPath(respJson, "aaa.mp3", "page", "0", "name")
	th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
	th.assertEqualsJsonPath(respJson, "ccc", "page", "2", "name")
	th.assertEqualsJsonPath(respJson, "ddd", "page", "3", "name")
	th.assertEqualsJsonPath(respJson, "eee", "page", "4", "name")
}
