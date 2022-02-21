package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadMetadata(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return
	}

	f := addFile(newFileRequest)

	var response uploadMetadataResponse
	response.Id = f.Id
	response.UploadUrl = f.Url

	c.IndentedJSON(http.StatusCreated, response)
}

func listFiles(c *gin.Context) {
	var page []listFileRecord
	for _, f := range filesMemStorage {
		page = append(page, listFileRecordOf(f))
	}
	c.IndentedJSON(http.StatusOK, listFilesResponse{
		Page:  page,
		Total: uint(len(filesMemStorage)),
	})
}

func main() {
	router := gin.Default()
	router.POST("/api/file/upload", uploadMetadata)
	router.GET("/api/file", listFiles)

	router.Run("localhost:8080")
}
