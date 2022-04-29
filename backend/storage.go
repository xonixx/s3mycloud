package main

type Storage interface {
	addFile(request uploadMetadataRequest) (file, error)
	listFiles(listQuery listFilesQueryRequest) (listFilesResponse, error)
	removeFile(id string) error
	assignTags(id string, tags []string) error
	removeTags(id string, tags []string) error
	cleanStorage() error
}

type file struct {
	Id      string
	Name    string
	Size    uint
	Tags    map[string]bool
	Url     string // S3 URL todo
	Created int64
}

func (f *file) tags() []string {
	tags := make([]string, 0)
	for t := range f.Tags {
		tags = append(tags, t)
	}
	return tags
}