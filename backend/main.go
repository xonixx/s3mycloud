package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func uploadMetadataHandler(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return
	}

	// TODO allow explicit size=0 case
	f := addFile(newFileRequest)

	var response uploadMetadataResponse
	response.Id = f.Id
	response.UploadUrl = f.Url

	c.IndentedJSON(http.StatusCreated, response)
}

func listFilesHandler(c *gin.Context) {
	var listQuery listFilesQueryRequest
	if err := c.BindQuery(&listQuery); err != nil {
		return
	}
	if listQuery.PageSize == 0 {
		listQuery.PageSize = 10
	}
	if len(listQuery.Sort) == 0 {
		listQuery.Sort = []string{"uploaded", "desc"}
	} else {
		listQuery.Sort = strings.Split(listQuery.Sort[0], ",")
	}
	if len(listQuery.Tags) > 0 {
		listQuery.Tags = strings.Split(listQuery.Tags[0], ",")
	}
	//fmt.Println("listQuery:", listQuery)
	c.IndentedJSON(http.StatusOK, listFiles(listQuery))
}

func deleteFileHandler(c *gin.Context) {
	id := c.Param("id")
	if e := removeFile(id); e != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(e))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func assignTagsHandler(c *gin.Context) {
	var tags listOfTags
	if err := c.BindJSON(&tags); err != nil {
		return
	}
	id := c.Param("id")
	if err := assignTags(id, tags); err != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(err))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func removeTagsHandler(c *gin.Context) {
	var tags listOfTags
	if err := c.BindJSON(&tags); err != nil {
		return
	}
	id := c.Param("id")
	if err := removeTags(id, tags); err != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(err))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func main() {
	router := gin.Default()

	// TODO use https://github.com/gin-gonic/gin#grouping-routes
	router.POST("/api/file/upload", uploadMetadataHandler)
	router.GET("/api/file", listFilesHandler)
	router.DELETE("/api/file/:id", deleteFileHandler)
	router.POST("/api/file/:id/tags", assignTagsHandler)
	router.DELETE("/api/file/:id/tags", removeTagsHandler)

	router.Run("localhost:8080")
}
