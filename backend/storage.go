package main

type Storage interface {
	addFile(request uploadMetadataRequest) file
	listFiles(listQuery listFilesQueryRequest) listFilesResponse
	removeFile(id string) error
	assignTags(id string, tags []string) error
	removeTags(id string, tags []string) error
	cleanStorage()
}

type file struct {
	Id      string
	Name    string
	Size    uint
	Tags    map[string]bool
	Url     string // S3 URL todo
	Created int64
}
