package storage

import "time"

type Storage interface {
	AddFile(fileData FileData) (StoredFile, error)
	ListFiles(listQuery ListFilesQuery) (ListFilesResult, error)
	GetFile(id string) (StoredFile, error)
	RemoveFile(id string) error
	AssignTags(id string, tags []string) error
	RemoveTags(id string, tags []string) error
	CleanStorage() error
}

type FileData struct {
	Name       string
	Size       uint
	Tags       []string
	Created    *time.Time
	ExternalId string
}

type StoredFile struct {
	Id         string
	Name       string
	Size       uint
	Tags       map[string]bool
	ExternalId string // S3 object key in case of S3
	Created    int64
}

func (f *StoredFile) GetTags() []string {
	tags := make([]string, 0)
	for t := range f.Tags {
		tags = append(tags, t)
	}
	return tags
}

type ListFilesQuery struct {
	Name     string
	Page     uint
	PageSize uint
	Tags     []string
	Sort     []string
}

type ListFilesResult struct {
	Page  []StoredFile
	Total uint
}
