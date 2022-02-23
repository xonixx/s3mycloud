package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
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
	t         *testing.T
	ts        *httptest.Server
	fileJsons []M
}

func (th testHelper) getExpectStatus(path string, expectedStatus int) M {
	if resp, err := http.Get(fmt.Sprintf("%s/%s", th.ts.URL, path)); err != nil {
		th.t.FailNow()
	} else {
		th.assertEquals(expectedStatus, resp.StatusCode)
		return tryGetBodyJson(th, resp)
	}
	return nil
}
func (th testHelper) deleteExpectStatus(path string, bodyJson interface{}, expectedStatus int) M {
	var body io.Reader
	if bodyJson != nil {
		body = toJson(bodyJson)
	}
	if req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", th.ts.URL, path), body); err != nil {
		th.t.FailNow()
	} else {
		if resp, err1 := http.DefaultClient.Do(req); err1 != nil {
			th.t.FailNow()
		} else {
			th.assertEquals(expectedStatus, resp.StatusCode)
			return tryGetBodyJson(th, resp)
		}
	}
	return nil
}
func (th testHelper) postExpectStatus(path string, bodyJson interface{}, expectedStatus int) M {
	if resp, err := http.Post(fmt.Sprintf("%s/%s", th.ts.URL, path), "application/json", toJson(bodyJson)); err != nil {
		th.t.FailNow()
	} else {
		th.assertEquals(expectedStatus, resp.StatusCode)
		return tryGetBodyJson(th, resp)
	}
	return nil
}
func tryGetBodyJson(th testHelper, resp *http.Response) M {
	val, ok := resp.Header["Content-Type"]

	if ok && val[0] == "application/json; charset=utf-8" {
		return readJsonAsMap(th.t, resp)
	}

	return nil
}

func (th testHelper) assertEqualsJsonPath(json M, expectedVal interface{}, path ...string) {
	val := query(json, path...)
	if _, ok := expectedVal.(int); ok {
		if floatVal, ok1 := val.(float64); ok1 {
			val = int(floatVal)
		}
	}
	th.assertEquals(expectedVal, val)
}

func (th testHelper) assertEquals(expected interface{}, actual interface{}) {
	if expected != actual {
		debug.PrintStack()
		th.t.Fatalf("expected %v != actual %v", expected, actual)
	}
}

func withSampleFiles(t *testing.T, testLogic func(th testHelper)) {
	withTestHelper(t, func(th testHelper) {
		th.setupFiles()
		testLogic(th)
	})
}

func withTestHelper(t *testing.T, testLogic func(th testHelper)) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()
	th := testHelper{t, ts, nil}

	cleanupMemStorage()

	testLogic(th)
}

const existingTag = "document"
const nonExistingTag = "tagDoesNotExist"

func (th *testHelper) setupFiles() {
	th.fileJsons = nil
	for _, f := range []M{
		{
			"name": "aaa.mp3",
			"size": 1000_000,
			"tags": []string{"music", "pop"},
		},
		{
			"name": "bbb.txt",
			"size": 100,
			"tags": []string{"text", existingTag},
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
			"name": "eee.txt",
			"size": 500,
			"tags": []string{"text"},
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/api/file/upload", th.ts.URL), "application/json", toJson(f))

		//fmt.Println("resp JSON:", readJsonAsMap(t, resp))

		respJson := checkAndReadRespJson(th.t, resp, err, http.StatusCreated)
		th.fileJsons = append(th.fileJsons, respJson)
		getJsonField(th.t, respJson, "id")
		getJsonField(th.t, respJson, "uploadUrl")
	}
}

func (th testHelper) existingId() string {
	return query(th.fileJsons[1], "id").(string)
}
func (th testHelper) nonExistingId() string {
	return "doesNotExist"
}

func TestUploadFileSuccess(t *testing.T) {
	cleanupMemStorage()
	ts := httptest.NewServer(setupServer())
	defer ts.Close()
	th := testHelper{t, ts, nil}

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

	th.assertEqualsJsonPath(respJson, 1, "total")

	th.assertEqualsJsonPath(respJson, "file.txt", "page", "0", "name")
	th.assertEqualsJsonPath(respJson, id, "page", "0", "id")

}

func TestListAllDefaults(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEqualsJsonPath(respJson, "eee.txt", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "ddd", "page", "1", "name")
		th.assertEqualsJsonPath(respJson, "ccc", "page", "2", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "3", "name")
		th.assertEqualsJsonPath(respJson, "aaa.mp3", "page", "4", "name")
	})
}

func TestListOlderFirst(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?sort=uploaded,asc", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEqualsJsonPath(respJson, "aaa.mp3", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
		th.assertEqualsJsonPath(respJson, "ccc", "page", "2", "name")
		th.assertEqualsJsonPath(respJson, "ddd", "page", "3", "name")
		th.assertEqualsJsonPath(respJson, "eee.txt", "page", "4", "name")
	})
}

func TestListDefaultPageSize(t *testing.T) { // mock default pageSize
	withSampleFiles(t, func(th testHelper) {
		keepDefaultPageSize := DefaultPageSize
		DefaultPageSize = 2
		defer func() {
			DefaultPageSize = keepDefaultPageSize
		}()

		respJson := th.getExpectStatus("api/file", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEquals(int(DefaultPageSize), len(query(respJson, "page").([]interface{})))

		th.assertEqualsJsonPath(respJson, "eee.txt", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "ddd", "page", "1", "name")
	})
}
func TestListPaging(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?page=1&pageSize=2", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEquals(2, len(query(respJson, "page").([]interface{})))
		th.assertEqualsJsonPath(respJson, "ccc", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
	})
}
func TestListPagingSorting(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?page=1&pageSize=2&sort=size,desc", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEquals(2, len(query(respJson, "page").([]interface{})))
		th.assertEqualsJsonPath(respJson, "ddd", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
	})
}
func TestListFilterSingleTag(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?tags=text", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 2, "total")

		th.assertEquals(2, len(query(respJson, "page").([]interface{})))
		th.assertEqualsJsonPath(respJson, "eee.txt", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
	})
}
func TestListFilterMultipleTagsWithAndLogic(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?tags=text,"+existingTag, http.StatusOK)

		th.assertEqualsJsonPath(respJson, 1, "total")

		th.assertEquals(1, len(query(respJson, "page").([]interface{})))
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "0", "name")
	})
}
func TestListFilterMultipleTagsWithAndLogicNotFound(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?tags=text,"+nonExistingTag, http.StatusOK)

		th.assertEqualsJsonPath(respJson, 0, "total")

		th.assertEquals(0, len(query(respJson, "page").([]interface{})))
	})
}
func TestListFilterName(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?name=txt", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 2, "total")

		th.assertEquals(2, len(query(respJson, "page").([]interface{})))
		th.assertEqualsJsonPath(respJson, "eee.txt", "page", "0", "name")
		th.assertEqualsJsonPath(respJson, "bbb.txt", "page", "1", "name")
	})
}
func TestListFilterNameNotFound(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?name=absentName", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 0, "total")

		th.assertEquals(0, len(query(respJson, "page").([]interface{})))
	})
}
func TestListComplex(t *testing.T) {
}
func TestListPageNegativeProduces400(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		th.getExpectStatus("api/file?page=-7", http.StatusBadRequest)
	})
}
func TestListPageNotANumberProduces400(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		th.getExpectStatus("api/file?page=P", http.StatusBadRequest)
	})
}
func TestListPageVeryLargeProducesEmptyPage(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.getExpectStatus("api/file?page=1000000", http.StatusOK)

		th.assertEqualsJsonPath(respJson, 5, "total")

		th.assertEquals(0, len(query(respJson, "page").([]interface{})))
	})
}
func TestListPageSizeNegativeProduces400(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		th.getExpectStatus("api/file?pageSize=-7", http.StatusBadRequest)
	})
}
func TestUploadEmptyNameProduces400(t *testing.T) {
	withTestHelper(t, func(th testHelper) {
		th.postExpectStatus("api/file/upload", M{
			"name": "",
			"size": 100,
		}, http.StatusBadRequest)
	})
}
func TestUploadOmittedNameProduces400(t *testing.T) {
	withTestHelper(t, func(th testHelper) {
		th.postExpectStatus("api/file/upload", M{
			"size": 100,
		}, http.StatusBadRequest)
	})
}
func TestUploadOmittedSizeProduces400(t *testing.T) {
	withTestHelper(t, func(th testHelper) {
		th.postExpectStatus("api/file/upload", M{
			"name": "name.txt",
		}, http.StatusBadRequest)
	})
}
func TestUploadZeroSizeIsOk(t *testing.T) {
	withTestHelper(t, func(th testHelper) {
		th.postExpectStatus("api/file/upload", M{
			"name": "name.txt",
			"size": 0,
		}, http.StatusCreated)
	})
}
func TestUploadNegativeSizeProduces400(t *testing.T) {
	withTestHelper(t, func(th testHelper) {
		th.postExpectStatus("api/file/upload", M{
			"name": "name.txt",
			"size": -100,
		}, http.StatusBadRequest)
	})
}
func TestDeleteExistingOk(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.deleteExpectStatus(fmt.Sprintf("api/file/%s", th.existingId()), nil, http.StatusOK)
		th.assertEqualsJsonPath(respJson, true, "success")
	})
}
func TestDeleteWrongIdProduces404(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.deleteExpectStatus(fmt.Sprintf("api/file/%s", th.nonExistingId()), nil, http.StatusNotFound)
		th.assertEqualsJsonPath(respJson, "file not found", "error")
	})
}
func TestAssignTagOk(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.postExpectStatus(fmt.Sprintf("api/file/%s/tags", th.existingId()), []string{"tag2", "tag3"}, http.StatusOK)
		th.assertEqualsJsonPath(respJson, true, "success")
	})
}
func TestAssignTagWrongFileProduces404(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.postExpectStatus(fmt.Sprintf("api/file/%s/tags", th.nonExistingId()), []string{"tag2", "tag3"}, http.StatusNotFound)
		th.assertEqualsJsonPath(respJson, "file not found", "error")
	})
}
func TestRemoveTagOk(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.deleteExpectStatus(fmt.Sprintf("api/file/%s/tags", th.existingId()), []string{existingTag}, http.StatusOK)
		th.assertEqualsJsonPath(respJson, true, "success")
	})
}
func TestRemoveTagWrongFileProduces404(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.deleteExpectStatus(fmt.Sprintf("api/file/%s/tags", th.nonExistingId()), []string{existingTag}, http.StatusNotFound)
		th.assertEqualsJsonPath(respJson, "file not found", "error")
	})
}
func TestRemoveWrongTagProduces404(t *testing.T) {
	withSampleFiles(t, func(th testHelper) {
		respJson := th.deleteExpectStatus(fmt.Sprintf("api/file/%s/tags", th.existingId()), []string{nonExistingTag}, http.StatusNotFound)
		th.assertEqualsJsonPath(respJson, "tag not found", "error")
	})
}
