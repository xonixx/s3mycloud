package main

import (
	"errors"
	"strconv"
	"time"
)

type file struct {
	Id      string
	Name    string
	Size    uint
	Tags    []string
	Url     string // S3 URL todo
	Created int64
}

var filesMemStorage []file
var globalId uint64

// @returns storage-level file object
func addFile(request uploadMetadataRequest) file {
	var f file
	f.Name = request.Name
	f.Size = request.Size
	f.Tags = request.Tags
	f.Url = "https://S3/todo"

	globalId += 1
	f.Id = strconv.FormatUint(globalId, 10)
	f.Created = time.Now().UnixNano()

	filesMemStorage = append(filesMemStorage, f)

	return f
}

func removeFile(id string) error {
	for i, f := range filesMemStorage {
		if id == f.Id {
			filesMemStorage = append(filesMemStorage[:i], filesMemStorage[i+1:]...)
			return nil
		}
	}
	return errors.New("file not found")
}
