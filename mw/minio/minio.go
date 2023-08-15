package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
)

var (
	Client *minio.Client
	err    error
)

// MakeBucket creates a new bucket with a specific name.
func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully create bucket %v\n", bucketName)
	}
}

// PutToBucket puts an object from a file into a bucket by *multipart.FileHeader.
func PutToBucket(ctx context.Context, bucketName string, file *multipart.FileHeader) (info minio.UploadInfo, err error) {
	fileObj, _ := file.Open()
	info, err = Client.PutObject(ctx, bucketName, file.Filename, fileObj, file.Size, minio.PutObjectOptions{})
	fileObj.Close()
	return info, err
}

// GetObjURL get the original link of the file in minio
func GetObjURL(ctx context.Context, bucketName string, filename string) (u *url.URL, err error) {
	expire := time.Hour * 24
	reqParams := make(url.Values)
	u, err = Client.PresignedGetObject(ctx, bucketName, filename, expire, reqParams)
	return u, err
}

// PutToBucketByBuf puts an object from a file into a bucket by *bytes.Buffer.
func PutToBucketByBuf(ctx context.Context, bucketName string, filename string, buf *bytes.Buffer) (info minio.UploadInfo, err error) {
	info, err = Client.PutObject(ctx, bucketName, filename, buf, int64(buf.Len()), minio.PutObjectOptions{})
	return info, err
}

// PutToBucketByFilePath put the file into bucket by filepath
func PutToBucketByFilePath(ctx context.Context, bucketName string, filename string, filepath string) (info minio.UploadInfo, err error) {
	info, err = Client.FPutObject(ctx, bucketName, filename, filepath, minio.PutObjectOptions{})
	return info, err
}

func Init() {
	ctx := context.Background()
	Client, err = minio.New(constants.MinioDefaultEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(constants.MinioAccessKeyID, constants.MinioSecrectAccessKey, ""),
		Secure: constants.MiniouseSSL,
	})
	if err != nil {
		log.Fatalln("minio connection error: ", err)
	}

	log.Printf("%#v\n", Client)

	MakeBucket(ctx, constants.MinioVideoBucketName)
	MakeBucket(ctx, constants.MinioImageBucketName)
}
