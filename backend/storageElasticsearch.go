package main

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

const INDEX = "file"

type storageElasticsearch struct {
	filesMemStorage []file
	globalId        uint64
	esClient        *elasticsearch.Client
}

type esFile struct {
	Id     string       `json:"_id"`
	Source esFileSource `json:"_source"`
}
type esFileSource struct {
	Name    string   `json:"name"`
	Size    uint     `json:"size"`
	Tags    []string `json:"tags"`
	Url     string   `json:"url"` // S3 URL todo
	Created int64    `json:"created"`
}
type esSearchResult struct {
	Took     uint `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value uint `json:"value"`
		} `json:"total"`
		Hits []esFile `json:"hits"`
	} `json:"hits"`
}

func NewElasticsearchStorage() Storage {
	if esClient, err := elasticsearch.NewDefaultClient(); err != nil {
		log.Fatalf("Can't connect ES: %v", err)
		return nil
	} else {
		//log.Println(esClient)
		//log.Println(esClient.Info())
		return &storageElasticsearch{esClient: esClient}
	}
}

func (s *storageElasticsearch) cleanStorage() error {
	s.filesMemStorage = nil

	//_, err := s.esClient.Indices.Delete([]string{INDEX})
	//if err != nil {
	//	log.Fatalf("Unable to delete index: %v", err)
	//}

	// This is faster
	resp, err := s.esClient.DeleteByQuery([]string{INDEX}, strings.NewReader(`{"query":{"match_all":{}}}`),
		s.esClient.DeleteByQuery.WithRefresh(true))
	if err1 := checkError(resp, err); err1 != nil {
		return err1
	}
	return nil
}

func toEsFileSource(f file) esFileSource {
	return esFileSource{
		//Id:      f.Id,
		Name:    f.Name,
		Size:    f.Size,
		Tags:    stringKeys(f.Tags),
		Url:     f.Url,
		Created: f.Created,
	}
}

func fromEsFile(ef esFile) file {
	source := ef.Source
	f := file{
		Id:      ef.Id,
		Name:    source.Name,
		Size:    source.Size,
		Tags:    map[string]bool{},
		Url:     source.Url,
		Created: source.Created,
	}
	for _, tag := range source.Tags {
		f.Tags[tag] = true
	}
	return f
}

// @returns storage-level file object
func (s *storageElasticsearch) addFile(request uploadMetadataRequest) (file, error) {
	var f file
	f.Name = request.Name
	f.Size = *request.Size
	f.Tags = map[string]bool{}
	for _, tag := range request.Tags {
		f.Tags[tag] = true
	}
	f.Url = "https://S3/todo"

	s.globalId += 1
	//f.Id = strconv.FormatUint(s.globalId, 10)
	f.Created = time.Now().UnixNano()

	indexResp, err := s.esClient.Index(INDEX, toJson(toEsFileSource(f)),
		s.esClient.Index.WithRefresh("true"))
	if err1 := checkError(indexResp, err); err1 != nil {
		return file{}, err1
	}
	log.Println("resp: ", indexResp)
	var ef esFile
	parseJsonTyped(indexResp.Body, &ef)
	id := ef.Id
	f.Id = id
	s.filesMemStorage = append(s.filesMemStorage, f)
	log.Println("ID: ", id)

	return f, nil
}

func (s *storageElasticsearch) findFile(id string) (int, file) {
	for i, f := range s.filesMemStorage {
		if id == f.Id {
			return i, f
		}
	}
	return -1, file{}
}

func (s *storageElasticsearch) findFileEs(id string) (*file, error) {
	resp, err := s.esClient.Get(INDEX, id)
	if resp.StatusCode == 404 {
		return nil, nil
	}
	if err1 := checkError(resp, err); err1 != nil {
		return nil, err1
	}
	var ef esFile
	parseJsonTyped(resp.Body, &ef)
	f := fromEsFile(ef)
	return &f, nil
}

func (s *storageElasticsearch) listFiles(listQuery listFilesQueryRequest) (listFilesResponse, error) {
	by := listQuery.Sort[0]
	desc := len(listQuery.Sort) > 1 && listQuery.Sort[1] == "desc"
	sort := by
	if by == "uploaded" {
		sort = "created"
	} //else if by == "name" {
	//sort = "name.keyword"
	//}
	order := "asc"
	if desc {
		order = "desc"
	}
	searchBody := M{
		"from": listQuery.Page * listQuery.PageSize,
		"size": listQuery.PageSize,
		"sort": M{sort: order},
	}
	conditions := make([]M, 0)
	if listQuery.Name != "" {
		conditions = append(conditions,
			M{
				"query_string": M{
					"default_field": "name",
					"query":         "*" + listQuery.Name + "*",
				},
			})
	}
	for _, tag := range listQuery.Tags {
		conditions = append(conditions, M{
			"term": M{"tags": tag},
		})
	}
	searchBody["query"] = M{
		"bool": M{
			"filter": conditions,
		},
	}
	log.Println("searchBody", searchBody)
	searchResp, err := s.esClient.Search(s.esClient.Search.WithIndex(INDEX),
		s.esClient.Search.WithBody(toJson(searchBody)))
	if err1 := checkError(searchResp, err); err1 != nil {
		return listFilesResponse{}, err1
	}
	log.Println("searchResp=", searchResp)
	var searchRes esSearchResult
	parseJsonTyped(searchResp.Body, &searchRes)

	page := make([]listFileRecord, 0)
	for _, ef := range searchRes.Hits.Hits {
		page = append(page, listFileRecordOf(fromEsFile(ef)))
	}

	return listFilesResponse{
		Total: searchRes.Hits.Total.Value,
		Page:  page,
	}, nil
}

func checkError(r *esapi.Response, err error) error {
	if err != nil {
		return err
	}
	if r.IsError() {
		return errors.New(fmt.Sprintf("%v", r))
	}
	return nil
}

func (s *storageElasticsearch) removeFile(id string) error {
	resp, err := s.esClient.Delete(INDEX, id)
	if err1 := checkError(resp, err); err1 != nil {
		return err1
	}

	if i, _ := s.findFile(id); i >= 0 {
		s.filesMemStorage = append(s.filesMemStorage[:i], s.filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) assignTags(id string, tags []string) error {
	f, err := s.findFileEs(id)
	if err != nil {
		return err
	}
	if f != nil {
		for _, t := range tags {
			f.Tags[t] = true
		}
		// TODO update file in ES
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) removeTags(id string, tags []string) error {
	f, err := s.findFileEs(id)
	if err != nil {
		return err
	}
	if f != nil {
		for _, t := range tags {
			if !f.Tags[t] {
				return errors.New("tag not found")
			}
		}
		for _, t := range tags {
			delete(f.Tags, t)
		}
		// TODO update file in ES
		return nil
	} else {
		return errors.New("file not found")
	}
}
