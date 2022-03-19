package main

type Storage interface {
	addFile(request uploadMetadataRequest) file
	listFiles(listQuery listFilesQueryRequest) listFilesResponse
	removeFile(id string) error
	assignTags(id string, tags []string) error
	removeTags(id string, tags []string) error
	cleanStorage()
}
