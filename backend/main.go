package main

import (
	"net/http"
	"s3mycloud/storage"
	. "s3mycloud/util"
	"strings"

	"github.com/gin-gonic/gin"
)

var s = createStorage()

func createStorage() storage.Storage {
	//return storage.NewMemStorage()
	return storage.NewElasticsearchStorage()
}

func uploadMetadataHandler(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return
	}

	f, err := s.AddFile(storage.FileData{
		Name: newFileRequest.Name,
		Size: *newFileRequest.Size,
		Tags: newFileRequest.Tags,
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

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
	files, err := s.ListFiles(storage.ListFilesQuery{
		Name:     listQuery.Name,
		Page:     listQuery.Page,
		PageSize: listQuery.PageSize,
		Tags:     listQuery.Tags,
		Sort:     listQuery.Sort,
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	} else {
		c.IndentedJSON(http.StatusOK, listFilesResponse{
			Page: Map(files.Page, func(f storage.StoredFile) listFileRecord {
				return listFileRecord{
					Id:       f.Id,
					Name:     f.Name,
					Size:     f.Size,
					Tags:     f.GetTags(),
					Url:      f.Url,
					Uploaded: f.Created,
				}
			}),
			Total: files.Total,
		})
	}
}

func deleteFileHandler(c *gin.Context) {
	id := c.Param("id")
	if e := s.RemoveFile(id); e != nil {
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
	if err := s.AssignTags(id, tags); err != nil {
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
	if err := s.RemoveTags(id, tags); err != nil {
		c.IndentedJSON(http.StatusNotFound, errorResponseOf(err))
		return
	}
	c.IndentedJSON(http.StatusOK, success)
}

func main() {
	err := setupServer().Run("localhost:8080")
	if err != nil {
		panic(err)
	}
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
