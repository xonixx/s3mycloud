package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"s3mycloud/util"
	"strings"
	"time"
)

type Config struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
}

type Connection struct {
	config Config
	client s3.Client
}

func Connect(conf Config) (Connection, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(conf.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: conf.Endpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, "")))
	if err != nil {
		//log.Fatal(err)
		return Connection{}, err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return Connection{client: *client, config: conf}, nil
}

func (conn *Connection) List() ([]File, error) {
	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := conn.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(conn.config.Bucket),
	})
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	//log.Println("first page results:")
	//for _, object := range output.Contents {
	//	log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	//}
	return util.Map(output.Contents, func(o types.Object) File {
		return s3FileOf(o)
	}), nil
}

func (conn *Connection) MakePreSignedUrl(key string) (string, error) {
	presignClient := s3.NewPresignClient(&conn.client)
	input := &s3.GetObjectInput{
		Bucket: &conn.config.Bucket,
		Key:    &key,
	}
	r, err := presignClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(30*time.Minute))
	if err != nil {
		return "", err
	}
	return r.URL, nil
}

type File struct {
	Key          string
	Size         int64
	LastModified time.Time
}

func s3FileOf(o types.Object) File {
	return File{
		Key:          *o.Key,
		Size:         o.Size,
		LastModified: *o.LastModified,
	}
}

func (s3f File) Path() string {
	parts := strings.Split(s3f.Key, "/")
	pathParts := parts[:len(parts)-1]
	return strings.Join(pathParts, "/")
}
func (s3f File) Name() string {
	parts := strings.Split(s3f.Key, "/")
	return parts[len(parts)-1]
}
