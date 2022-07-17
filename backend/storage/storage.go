package storage

type Storage interface {
	AddFile(request FileData) (StoredFile, error)
	ListFiles(listQuery ListFilesQuery) (ListFilesResult, error)
	RemoveFile(id string) error
	AssignTags(id string, tags []string) error
	RemoveTags(id string, tags []string) error
	CleanStorage() error
}

type FileData struct {
	Name string
	Size uint
	Tags []string
}

type StoredFile struct {
	Id      string
	Name    string
	Size    uint
	Tags    map[string]bool
	Url     string // S3 URL todo
	Created int64
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
