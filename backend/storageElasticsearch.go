package main

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"sort"
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
	Id      string
	Name    string
	Size    uint
	Tags    []string
	Url     string // S3 URL todo
	Created int64
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

	_, err := s.esClient.Indices.Delete([]string{INDEX})
	if err != nil {
		log.Fatalf("Unable to delete index: %v", err)
	}
}

func toEsFile(f file) esFile {
	return esFile{
		Id:      f.Id,
		Name:    f.Name,
		Size:    f.Size,
		Tags:    stringKeys(f.Tags),
		Url:     f.Url,
		Created: f.Created,
	}
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

	indexResp, err := s.esClient.Index(INDEX, toJson(toEsFile(f)))
	if err != nil {
		log.Fatalf("Unable to index: %v", err)
	}
	log.Println("resp: ", indexResp)
	id := parseJson(indexResp.Body)["_id"].(string)
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

func (s *storageElasticsearch) listFiles(listQuery listFilesQueryRequest) listFilesResponse {
	var matched []file
	for _, f := range s.filesMemStorage {
		var matchedName, matchedTags bool
		matchedName = true
		matchedTags = true
		if listQuery.Name != "" {
			matchedName = strings.Contains(f.Name, listQuery.Name)
		}
		if len(listQuery.Tags) > 0 {
			for _, tag := range listQuery.Tags {
				if !f.Tags[tag] {
					matchedTags = false
					break
				}
			}
		}
		if matchedName && matchedTags {
			matched = append(matched, f)
		}
	}
	total := len(matched)
	var page = make([]listFileRecord, 0)
	var from, to int
	from = int(listQuery.Page * listQuery.PageSize)
	to = from + int(listQuery.PageSize)
	if to > total {
		to = total
	}
	if len(listQuery.Sort) > 0 {
		by := listQuery.Sort[0]
		desc := len(listQuery.Sort) > 1 && listQuery.Sort[1] == "desc"
		sort.SliceStable(matched, func(i, j int) bool {
			if desc {
				i, j = j, i
			}
			f1 := matched[i]
			f2 := matched[j]
			if by == "name" {
				return f1.Name < f2.Name
			} else if by == "size" {
				return f1.Size < f2.Size
			} else if by == "uploaded" {
				return f1.Created < f2.Created
			} else {
				panic("Wrong sort by: " + by)
			}
		})
	}
	if from <= total {
		for _, f := range matched[from:to] {
			page = append(page, listFileRecordOf(f))
		}
	}
	return listFilesResponse{
		Page:  page,
		Total: uint(total),
	}
}

func (s *storageElasticsearch) removeFile(id string) error {
	if i, _ := s.findFile(id); i >= 0 {
		s.filesMemStorage = append(s.filesMemStorage[:i], s.filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) assignTags(id string, tags []string) error {
	if i, f := s.findFile(id); i >= 0 {
		for _, t := range tags {
			f.Tags[t] = true
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (s *storageElasticsearch) removeTags(id string, tags []string) error {
	if i, f := s.findFile(id); i >= 0 {
		for _, t := range tags {
			if !f.Tags[t] {
				return errors.New("tag not found")
			}
		}
		for _, t := range tags {
			delete(f.Tags, t)
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}
