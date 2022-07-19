package storage

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type memoryStorage struct {
	filesMemStorage []StoredFile
	globalId        uint64
}

func NewMemStorage() Storage {
	return &memoryStorage{}
}

func (m *memoryStorage) CleanStorage() error {
	m.filesMemStorage = nil
	return nil
}

// @returns storage-level file object
func (m *memoryStorage) AddFile(fileData FileData) (StoredFile, error) {
	var f StoredFile
	f.Name = fileData.Name
	f.Size = fileData.Size
	f.Tags = map[string]bool{}
	for _, tag := range fileData.Tags {
		f.Tags[tag] = true
	}
	f.ExternalId = fileData.ExternalId

	m.globalId += 1
	f.Id = strconv.FormatUint(m.globalId, 10)

	if fileData.Created == nil {
		return f, errors.New("created should be set")
	}
	f.Created = fileData.Created.UnixMilli()

	m.filesMemStorage = append(m.filesMemStorage, f)

	return f, nil
}

func (m *memoryStorage) findFile(id string) (int, StoredFile) {
	for i, f := range m.filesMemStorage {
		if id == f.Id {
			return i, f
		}
	}
	return -1, StoredFile{}
}

func (m *memoryStorage) ListFiles(listQuery ListFilesQuery) (ListFilesResult, error) {
	var matched []StoredFile
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
	var page = make([]StoredFile, 0)
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
			page = append(page, f)
		}
	}
	return ListFilesResult{
		Page:  page,
		Total: uint(total),
	}, nil
}

func (m *memoryStorage) RemoveFile(id string) error {
	if i, _ := m.findFile(id); i >= 0 {
		m.filesMemStorage = append(m.filesMemStorage[:i], m.filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (m *memoryStorage) AssignTags(id string, tags []string) error {
	if i, f := m.findFile(id); i >= 0 {
		for _, t := range tags {
			f.Tags[t] = true
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}

func (m *memoryStorage) RemoveTags(id string, tags []string) error {
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
