package storage

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	. "s3mycloud/util"
	"strings"
)

const INDEX = "file"

type storageElasticsearch struct {
	esClient *elasticsearch.Client
}

type esFile struct {
	Id     string       `json:"_id"`
	Source esFileSource `json:"_source"`
}
type esFileSource struct {
	Name       string   `json:"name"`
	Size       uint     `json:"size"`
	Tags       []string `json:"tags"`
	ExternalId string   `json:"externalId"`
	Created    int64    `json:"created"`
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

func (s *storageElasticsearch) CleanStorage() error {
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

func toEsFileSource(f StoredFile) esFileSource {
	return esFileSource{
		//Id:      f.Id,
		Name:       f.Name,
		Size:       f.Size,
		Tags:       f.GetTags(),
		ExternalId: f.ExternalId,
		Created:    f.Created,
	}
}

func fromEsFile(ef esFile) StoredFile {
	source := ef.Source
	f := StoredFile{
		Id:         ef.Id,
		Name:       source.Name,
		Size:       source.Size,
		Tags:       map[string]bool{},
		ExternalId: source.ExternalId,
		Created:    source.Created,
	}
	for _, tag := range source.Tags {
		f.Tags[tag] = true
	}
	return f
}

// @returns storage-level file object
func (s *storageElasticsearch) AddFile(fileData FileData) (StoredFile, error) {
	var f StoredFile
	f.Name = fileData.Name
	f.Size = fileData.Size
	f.Tags = map[string]bool{}
	for _, tag := range fileData.Tags {
		f.Tags[tag] = true
	}
	f.ExternalId = fileData.ExternalId

	//f.Id = strconv.FormatUint(s.globalId, 10)
	if fileData.Created == nil {
		return f, errors.New("created should be set")
	}
	f.Created = fileData.Created.UnixMilli()

	indexResp, err := s.esClient.Index(INDEX, ToJson(toEsFileSource(f)),
		s.esClient.Index.WithRefresh("true"))
	if err1 := checkError(indexResp, err); err1 != nil {
		return StoredFile{}, err1
	}
	log.Println("resp: ", indexResp)
	var ef esFile
	ParseJsonTyped(indexResp.Body, &ef)
	id := ef.Id
	f.Id = id
	log.Println("ID: ", id)

	return f, nil
}

func (s *storageElasticsearch) findFileEs(id string) (*StoredFile, error) {
	resp, err := s.esClient.Get(INDEX, id)
	if resp.StatusCode == 404 {
		return nil, nil
	}
	if err1 := checkError(resp, err); err1 != nil {
		return nil, err1
	}
	var ef esFile
	ParseJsonTyped(resp.Body, &ef)
	f := fromEsFile(ef)
	return &f, nil
}

func (s *storageElasticsearch) ListFiles(listQuery ListFilesQuery) (ListFilesResult, error) {
	by := listQuery.Sort[0]
	desc := len(listQuery.Sort) > 1 && listQuery.Sort[1] == "desc"
	sort := by
	if by == "uploaded" {
		sort = "created"
	} else if by == "name" {
		sort = "name.keyword"
	}
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
		s.esClient.Search.WithBody(ToJson(searchBody)))
	if err1 := checkError(searchResp, err); err1 != nil {
		return ListFilesResult{}, err1
	}
	log.Println("searchResp=", searchResp)
	var searchRes esSearchResult
	ParseJsonTyped(searchResp.Body, &searchRes)

	page := make([]StoredFile, 0)
	for _, ef := range searchRes.Hits.Hits {
		page = append(page, fromEsFile(ef))
	}

	return ListFilesResult{
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

func (s *storageElasticsearch) RemoveFile(id string) error {
	resp, err := s.esClient.Delete(INDEX, id, s.esClient.Delete.WithRefresh("true"))
	if resp.StatusCode == 404 {
		return errors.New("file not found")
	}
	if err1 := checkError(resp, err); err1 != nil {
		return err1
	}
	return nil
}

func (s *storageElasticsearch) AssignTags(id string, tags []string) error {
	f, err := s.findFileEs(id)
	if err != nil {
		return err
	}
	if f != nil {
		for _, t := range tags {
			f.Tags[t] = true
		}
		err := s.updateTags(f)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) updateTags(f *StoredFile) error {
	resp, err := s.esClient.Update(INDEX, f.Id, ToJson(M{
		"doc": M{
			"tags": f.GetTags(),
		},
	}), s.esClient.Update.WithRefresh("true"))
	if err1 := checkError(resp, err); err1 != nil {
		return err1
	}
	return nil
}

func (s *storageElasticsearch) RemoveTags(id string, tags []string) error {
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
		err := s.updateTags(f)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}
