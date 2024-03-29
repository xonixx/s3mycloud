package main

import "s3mycloud/storage"

type uploadMetadataRequest struct {
	Name string   `json:"name" binding:"required"`
	Size *uint    `json:"size" binding:"required"` // using pointer to distinguish size=0 from omitted
	Tags []string `json:"tags"`
}

type uploadMetadataResponse struct {
	Id        string `json:"id,omitempty"`
	UploadUrl string `json:"uploadUrl,omitempty"`
	Error     string `json:"error,omitempty"`
}

type confirmUploadRequest struct {
	Id      string `json:"id" binding:"required"`
	Success *bool  `json:"success" binding:"required"` // using pointer to distinguish false from omitted
	Error   string `json:"error"`
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

type successResponse struct {
	Success bool `json:"success"`
}

type listFilesQueryRequest struct {
	Name     string   `form:"name"`
	Page     uint     `form:"page"`
	PageSize uint     `form:"pageSize"`
	Tags     []string `form:"tags"`
	Sort     []string `form:"sort"`
}

type listFileRecord struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Size     uint     `json:"size"`
	Tags     []string `json:"tags"`
	Url      string   `json:"url"`
	Uploaded int64    `json:"uploaded"`
}

type listFilesResponse struct {
	Page  []listFileRecord `json:"page"`
	Total uint             `json:"total"`
}

type listOfTags []string

func listFileRecordOf(f storage.StoredFile) listFileRecord {
	return listFileRecord{
		Id:       f.Id,
		Name:     f.Name,
		Size:     f.Size,
		Tags:     f.GetTags(),
		Uploaded: f.Created,
	}
}

func errorResponseOf(e error) errorResponse {
	return errorResponse{e.Error()}
}

var success = successResponse{true}
