package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadMetadata(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return // TODO 400
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

func deleteFile(c *gin.Context) {
	id := c.Param("id")
	if e := removeFile(id); e != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(e))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func apiAssignTags(c *gin.Context) {
	var tags listOfTags
	if err := c.BindJSON(&tags); err != nil {
		return // TODO 400
	}
	id := c.Param("id")
	if err := assignTags(id, tags); err != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(err))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func main() {
	router := gin.Default()
	router.POST("/api/file/upload", uploadMetadata)
	router.GET("/api/file", listFiles)
	router.DELETE("/api/file/:id", deleteFile)
	router.POST("/api/file/:id/tags", apiAssignTags)

	router.Run("localhost:8080")
}
