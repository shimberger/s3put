package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
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
	region      = kingpin.Flag("region", "The region to use.").Required().String()
	bucket      = kingpin.Flag("bucket", "The bucket to upload into.").Required().String()
	source      = kingpin.Arg("source", "The file to upload.").Required().String()
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

	// Read the file or stdin
	var file io.Reader = nil
	var err error = nil
	if *source == "-" {
		file = bufio.NewReader(os.Stdin)
	} else {
		file, err = os.Open(*source)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(*bucket),
		Key:         aws.String(*target),
		ContentType: contentType,
		Metadata: map[string]*string{
			"Content-Type": contentType,
		},
		ACL:  acl,
		Body: file,
	})

	// Check for errors
	if err != nil {
		fmt.Printf("failed to upload file, %v\n", err)
		return
	}

	fmt.Printf("Uploaded to %v\n", result.Location)
}
