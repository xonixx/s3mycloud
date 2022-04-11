package main

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

type memoryStorage struct {
	filesMemStorage []file
	globalId        uint64
}

func NewMemStorage() Storage {
	return &memoryStorage{}
}

func (m *memoryStorage) cleanStorage() error {
	m.filesMemStorage = nil
	return nil
}

// @returns storage-level file object
func (m *memoryStorage) addFile(request uploadMetadataRequest) (file, error) {
	var f file
	f.Name = request.Name
	f.Size = *request.Size
	f.Tags = map[string]bool{}
	for _, tag := range request.Tags {
		f.Tags[tag] = true
	}
	f.Url = "https://S3/todo"

	m.globalId += 1
	f.Id = strconv.FormatUint(m.globalId, 10)
	f.Created = time.Now().UnixNano()

	m.filesMemStorage = append(m.filesMemStorage, f)

	return f, nil
}

func (m *memoryStorage) findFile(id string) (int, file) {
	for i, f := range m.filesMemStorage {
		if id == f.Id {
			return i, f
		}
	}
	return -1, file{}
}

func (m *memoryStorage) listFiles(listQuery listFilesQueryRequest) (listFilesResponse, error) {
	var matched []file
	for _, f := range m.filesMemStorage {
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
	}, nil
}

func (m *memoryStorage) removeFile(id string) error {
	if i, _ := m.findFile(id); i >= 0 {
		m.filesMemStorage = append(m.filesMemStorage[:i], m.filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (m *memoryStorage) assignTags(id string, tags []string) error {
	if i, f := m.findFile(id); i >= 0 {
		for _, t := range tags {
			f.Tags[t] = true
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (m *memoryStorage) removeTags(id string, tags []string) error {
	if i, f := m.findFile(id); i >= 0 {
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
