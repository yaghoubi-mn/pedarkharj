package s3

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var accessKey string
var secretKey string
var apiUrlValue string
var bucketName string

var s3Client *s3.S3

func init() {
	// for _, s := range os.Environ() {
	// 	fmt.Println(s)
	// }

	accessKey = os.Getenv("S3_ACCESS_KEY")
	secretKey = os.Getenv("S3_SECRET_KEY")
	bucketName = os.Getenv("S3_BUCKET_NAME")
	apiUrlValue = os.Getenv("S3_API_URL_VALUE")
	fmt.Println(accessKey, secretKey, bucketName, apiUrlValue)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(apiUrlValue),
		Region:           aws.String("ir"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		slog.Error("cannot create s3 session: ", "errror", err)
	}

	s3Client = s3.New(newSession)

	if s3Client == nil {
		slog.Error("cannot create s3Client")
	}
}

// key exmaple: folder/name.format
func PutObject(key string, body io.ReadSeeker) error {
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:   body,
		Bucket: &bucketName,
		Key:    &key,
	})

	return err
}

// key example: https://domain.name/folder/name.format
func DeleteObject(key string) error {
	splited := strings.Split(key, apiUrlValue)
	if len(splited) > 1 {
		key = splited[1]
	}
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(key),
	})

	return err
}

// []string example: {"https://domain.name/folder/name.format", ... }
func GetListObjects(key string) ([]string, error) {
	resp, err := s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, err
	}

	list := make([]string, len(resp.Contents))
	for i, item := range resp.Contents {
		list[i] = path.Join(apiUrlValue, *item.Key)
	}

	return list, nil
}
