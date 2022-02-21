package main

type uploadMetadataRequest struct {
	Name string   `json:"name"`
	Size uint     `json:"size"`
	Tags []string `json:"tags"`
}

type uploadMetadataResponse struct {
	Id        string `json:"id,omitempty"`
	UploadUrl string `json:"uploadUrl,omitempty"`
	Error     string `json:"error,omitempty"`
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

type successResponse struct {
	Success bool `json:"success"`
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

func listFileRecordOf(f file) listFileRecord {
	return listFileRecord{
		Id:       f.Id,
		Name:     f.Name,
		Size:     f.Size,
		Url:      f.Url,
		Tags:     f.Tags,
		Uploaded: f.Created,
	}
}

func errorResponseOf(e error) errorResponse {
	return errorResponse{e.Error()}
}

var success = successResponse{true}
