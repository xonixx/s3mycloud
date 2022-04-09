package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var s = createStorage()

func createStorage() Storage {
	//return NewMemStorage()
	return NewElasticsearchStorage()
}

func uploadMetadataHandler(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return
	}

	f := s.addFile(newFileRequest)

	var response uploadMetadataResponse
	response.Id = f.Id
	response.UploadUrl = f.Url

	c.IndentedJSON(http.StatusCreated, response)
}

var DefaultPageSize uint = 10

func listFilesHandler(c *gin.Context) {
	var listQuery listFilesQueryRequest
	if err := c.BindQuery(&listQuery); err != nil {
		return
	}
	if listQuery.PageSize == 0 {
		listQuery.PageSize = DefaultPageSize
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
	c.IndentedJSON(http.StatusOK, s.listFiles(listQuery))
}

func deleteFileHandler(c *gin.Context) {
	id := c.Param("id")
	if e := s.removeFile(id); e != nil {
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
	if err := s.assignTags(id, tags); err != nil {
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
	if err := s.removeTags(id, tags); err != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(err))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func main() {
	setupServer().Run("localhost:8080")
}

func setupServer() *gin.Engine {
	router := gin.Default()

	// TODO use https://github.com/gin-gonic/gin#grouping-routes
	router.POST("/api/file/upload", uploadMetadataHandler)
	router.GET("/api/file", listFilesHandler)
	router.DELETE("/api/file/:id", deleteFileHandler)
	router.POST("/api/file/:id/tags", assignTagsHandler)
	router.DELETE("/api/file/:id/tags", removeTagsHandler)
	return router
}
