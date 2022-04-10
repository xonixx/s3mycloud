package main

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
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
	Name    string
	Size    uint
	Tags    []string
	Url     string // S3 URL todo
	Created int64
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

func (s *storageElasticsearch) cleanStorage() {
	s.filesMemStorage = nil

	//_, err := s.esClient.Indices.Delete([]string{INDEX})
	//if err != nil {
	//	log.Fatalf("Unable to delete index: %v", err)
	//}

	// This is faster
	_, err := s.esClient.DeleteByQuery([]string{INDEX}, strings.NewReader(`{"query":{"match_all":{}}}`),
		s.esClient.DeleteByQuery.WithRefresh(true))
	if err != nil {
		log.Fatalf("Unable to delete index: %v", err)
	}
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
func (s *storageElasticsearch) addFile(request uploadMetadataRequest) file {
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
	if err != nil {
		log.Fatalf("Unable to index: %v", err)
	}
	log.Println("resp: ", indexResp)
	var ef esFile
	parseJsonTyped(indexResp.Body, &ef)
	id := ef.Id
	f.Id = id
	s.filesMemStorage = append(s.filesMemStorage, f)
	log.Println("ID: ", id)

	// TODO check HasErrors

	return f
}

func (s *storageElasticsearch) findFile(id string) (int, file) {
	for i, f := range s.filesMemStorage {
		if id == f.Id {
			return i, f
		}
	}
	return -1, file{}
}

func (s *storageElasticsearch) findFileEs(id string) *file {
	resp, err := s.esClient.Get(INDEX, id)
	if err != nil {
		log.Fatalf("error get: %v", err)
	}
	if resp.StatusCode == 404 {
		return nil
	}
	var ef esFile
	parseJsonTyped(resp.Body, &ef)
	f := fromEsFile(ef)
	return &f
}

func (s *storageElasticsearch) listFiles(listQuery listFilesQueryRequest) listFilesResponse {
	searchResp, err := s.esClient.Search(s.esClient.Search.WithIndex(INDEX))
	if err != nil {
		log.Fatalf("Unable to search: %v", err)
	}
	log.Println("searchResp=", searchResp)
	var searchRes esSearchResult
	parseJsonTyped(searchResp.Body, &searchRes)

	var page []listFileRecord
	for _, ef := range searchRes.Hits.Hits {
		page = append(page, listFileRecordOf(fromEsFile(ef)))
	}

	return listFilesResponse{
		Total: searchRes.Hits.Total.Value,
		Page:  page,
	}
}

func (s *storageElasticsearch) removeFile(id string) error {
	_, err := s.esClient.Delete(INDEX, id)
	if err != nil {
		log.Printf("Unable to delete file: %v", err)
		return err
	}

	if i, _ := s.findFile(id); i >= 0 {
		s.filesMemStorage = append(s.filesMemStorage[:i], s.filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) assignTags(id string, tags []string) error {
	if f := s.findFileEs(id); f != nil {
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
	if f := s.findFileEs(id); f != nil {
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
