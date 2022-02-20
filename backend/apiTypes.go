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
		Url:      f.Url,
		Tags:     f.Tags,
		Uploaded: f.Created,
	}
}
