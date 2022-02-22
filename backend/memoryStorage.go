package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type file struct {
	Id      string
	Name    string
	Size    uint
	Tags    map[string]bool
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
	f.Tags = map[string]bool{}
	for _, tag := range request.Tags {
		f.Tags[tag] = true
	}
	f.Url = "https://S3/todo"

	globalId += 1
	f.Id = strconv.FormatUint(globalId, 10)
	f.Created = time.Now().UnixNano()

	filesMemStorage = append(filesMemStorage, f)

	return f
}

func findFile(id string) (int, file) {
	for i, f := range filesMemStorage {
		if id == f.Id {
			return i, f
		}
	}
	return -1, file{}
}

func listFiles(listQuery listFilesQueryRequest) listFilesResponse {
	var matched []file
	for _, f := range filesMemStorage {
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
	var page []listFileRecord
	var from, to int
	from = int(listQuery.Page * listQuery.PageSize)
	to = from + int(listQuery.PageSize)
	if to > total {
		to = total
	}
	for _, f := range matched[from:to] {
		page = append(page, listFileRecordOf(f))
	}
	return listFilesResponse{
		Page:  page,
		Total: uint(total),
	}
}

func removeFile(id string) error {
	if i, _ := findFile(id); i >= 0 {
		filesMemStorage = append(filesMemStorage[:i], filesMemStorage[i+1:]...)
		return nil
	} else {
		return errors.New("file not found")
	}
}

func assignTags(id string, tags []string) error {
	if i, f := findFile(id); i >= 0 {
		for _, t := range tags {
			f.Tags[t] = true
		}
		return nil
	} else {
		return errors.New("file not found")
	}
}

func removeTags(id string, tags []string) error {
	if i, f := findFile(id); i >= 0 {
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
