package main

import (
	"fmt"
	"net/http"
	"s3mycloud/s3"
	"s3mycloud/storage"
	. "s3mycloud/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var s = createStorage()

func createStorage() storage.Storage {
	return storage.NewMemStorage()
	//return storage.NewElasticsearchStorage()
}

func uploadMetadataHandler(c *gin.Context) {
	var newFileRequest uploadMetadataRequest

	if err := c.BindJSON(&newFileRequest); err != nil {
		return
	}

	now := time.Now()
	f, err := s.AddFile(storage.FileData{
		Name:    newFileRequest.Name,
		Size:    *newFileRequest.Size,
		Tags:    newFileRequest.Tags,
		Created: &now,
	})
	if err != nil {
		fmt.Println("err:", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	var response uploadMetadataResponse
	response.Id = f.Id
	response.UploadUrl, err = s3Connection.MakePreSignedPutUrl(newFileRequest.Name) // TODO should we add ID to limit the prob for collisions?
	if err != nil {
		fmt.Println("err:", err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

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
				return listFileRecordOf(f)
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

func downloadFileHandler(c *gin.Context) {
	id := c.Param("id")
	file, err := s.GetFile(id)
	// TODO we should distinguish 404 vs 500
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		preSigneUrl, err := s3Connection.MakePreSignedGetUrl(file.ExternalId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.Redirect(http.StatusFound, preSigneUrl)
	}
}

var myConfig Config
var s3Connection s3.Operations

func main() {
	{
		config, err := LoadConfig("../application-external.yml")
		if err != nil {
			panic(err)
		}
		myConfig = *config
		s3Connection, err = s3.Connect(s3.Config{
			Bucket:    myConfig.S3.Bucket,
			AccessKey: myConfig.S3.AccessKey,
			SecretKey: myConfig.S3.SecretKey,
			Endpoint:  myConfig.S3.Endpoint,
			Region:    myConfig.S3.Region,
		})
		if err != nil {
			panic(err)
		}
	}
	{
		err := addMockData()
		if err != nil {
			panic(err)
		}
	}
	err := setupServer().Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupServer() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// TODO use https://github.com/gin-gonic/gin#grouping-routes
	router.POST("/api/file/upload", uploadMetadataHandler)
	router.GET("/api/file", listFilesHandler)
	router.DELETE("/api/file/:id", deleteFileHandler)
	router.POST("/api/file/:id/tags", assignTagsHandler)
	router.DELETE("/api/file/:id/tags", removeTagsHandler)
	router.GET("/api/file/:id/dl", downloadFileHandler)
	return router
}

func addMockData() error {
	//fmt.Println(myConfig.S3.Bucket)
	s3Files, err := s3Connection.List()
	if err != nil {
		return err
	}
	for _, f := range s3Files {
		_, err := s.AddFile(storage.FileData{
			Name:       f.Name(),
			ExternalId: f.Key,
			Size:       uint(f.Size),
			Created:    &f.LastModified,
			Tags:       []string{f.Path()},
		})
		if err != nil {
			return err
		}
	}
	fmt.Printf("Got %d files from S3\n", len(s3Files))
	return nil
}

/*
func date(dateS string) int64 {
	d, err := time.Parse("2 Jan 2006", dateS)
	if err != nil {
		panic(err)
	}
	return d.UnixMilli()
}
func addMockData() {
	for i := 0; i < 2; i++ {
		s.AddFile(storage.FileData{Name: "Report for boss.xlsx", Size: 50000, Created: date("15 Mar 2016"), Tags: []string{"document", "work"}})
		s.AddFile(storage.FileData{Name: "Sing Now.mp3", Size: 2_500_000, Created: date("17 Apr 2019"), Tags: []string{"music", "pop"}})
		s.AddFile(storage.FileData{Name: "Test.txt", Size: 100, Created: date("1 Jan 2008"), Tags: []string{}})
		s.AddFile(storage.FileData{Name: "CV (John Doe).pdf", Size: 123_456, Created: date("9 Mar 2020"), Tags: []string{"work"}})
		s.AddFile(storage.FileData{Name: "Some veeeeeeery loooooooooooong naaaaaaaaame.ext", Size: 0, Created: date("31 Jan 2010"), Tags: []string{"test"}})
	}
}*/
