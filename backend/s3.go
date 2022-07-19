package main

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

// TODO this should be executed only once and cached
func getS3Client(conf Config) (s3.Client, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(conf.S3.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: conf.S3.Endpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.S3.AccessKey, conf.S3.SecretKey, "")))
	if err != nil {
		//log.Fatal(err)
		return s3.Client{}, err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	return *client, nil
}

func listS3(conf Config) ([]S3File, error) {
	client, err := getS3Client(conf)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(conf.S3.Bucket),
	})
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	//log.Println("first page results:")
	//for _, object := range output.Contents {
	//	log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	//}
	return util.Map(output.Contents, func(o types.Object) S3File {
		return s3FileOf(o)
	}), nil
}

func makePreSignedUrl(config Config, key string) (string, error) {
	client, err := getS3Client(config)
	if err != nil {
		return "", err
	}
	presignClient := s3.NewPresignClient(&client)
	input := &s3.GetObjectInput{
		Bucket: &config.S3.Bucket,
		Key:    &key,
	}
	r, err := presignClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(30*time.Minute))
	if err != nil {
		return "", err
	}
	return r.URL, nil
}

type S3File struct {
	Key          string
	Size         int64
	LastModified time.Time
}

func s3FileOf(o types.Object) S3File {
	return S3File{
		Key:          *o.Key,
		Size:         o.Size,
		LastModified: *o.LastModified,
	}
}

func (s3f S3File) Path() string {
	parts := strings.Split(s3f.Key, "/")
	pathParts := parts[:len(parts)-1]
	return strings.Join(pathParts, "/")
}
func (s3f S3File) Name() string {
	parts := strings.Split(s3f.Key, "/")
	return parts[len(parts)-1]
}
