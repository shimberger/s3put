package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	acl = kingpin.Flag("acl", "The canned ACL to use").Default("private").Enum(
		"private",
		"public-read",
		"public-read-write",
		"aws-exec-read",
		"authenticated-read",
		"bucket-owner-read",
		"bucket-owner-full-control",
		"log-delivery-write",
	)
	contentType = kingpin.Flag("content-type", "The content type to assign this object").Default("binary/octet-stream").String()
	partSize    = kingpin.Flag("part-size", "The part size in MB.").Default("64").Int64()
	region      = kingpin.Flag("region", "The region to use.").Required().String()
	bucket      = kingpin.Flag("bucket", "The bucket to upload into.").Required().String()
	source      = kingpin.Arg("source", "The file to upload.").Required().File()
	target      = kingpin.Arg("target", "The name of the file.").Required().String()
)

func main() {
	kingpin.Parse()

	// Create the AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*region),
	}))

	// Create an S3 uploader, this uploader allows multipart uploads
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(*bucket),
		Key:         aws.String(*target),
		ContentType: contentType,
		Metadata: map[string]*string{
			"Content-Type": contentType,
		},
		ACL:  acl,
		Body: *source,
	}, func(u *s3manager.Uploader) {
		u.PartSize = *partSize * 1024 * 1024
	})

	// Check for errors
	if err != nil {
		fmt.Printf("failed to upload file, %v\n", err)
		return
	}

	fmt.Printf("Uploaded to %v\n", result.Location)
}
