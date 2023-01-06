package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

var projectID = flag.String("project_id", "sample_project", "Project id for gcp storage")
var bucketName = flag.String("bucket_name", "sample_bucket", "Bucket name for gcp storage")
var uploadFilename = flag.String("file", "", "Filename for upload")

// var

func main() {
	ctx := context.Background()

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Errorf("Could not initiate the storage client. err: %v\n", err)
		return
	}

	// print buckets
	printBucketNames(ctx, client, *projectID)

	// upload file
	uploadFile(ctx, client, *bucketName, *uploadFilename)
}

func printBucketNames(ctx context.Context, client *storage.Client, projectID string) {

	it := client.Buckets(ctx, projectID)

	fmt.Printf("Bucket detail.\n")
	for {
		bucket, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			logrus.Errorf("Could not get bucket info correctly. err: %v\n", err)
			break
		}
		fmt.Printf("name: %s\n", bucket.Name)
	}

	return
}

// uploadFile uploads the file
func uploadFile(ctx context.Context, client *storage.Client, bucket string, filename string) error {
	// open file
	f, err := os.Open(filename)
	if err != nil {
		logrus.Errorf("Could not open the target file. err: %v", err)
		return err
	}
	defer f.Close()

	// tmp target diretory
	target := fmt.Sprintf("tmp/%s", filepath.Base(filename))

	// create a session
	wc := client.Bucket(bucket).Object(target).NewWriter(ctx)
	defer wc.Close()

	// upload the file
	if _, err = io.Copy(wc, f); err != nil {
		logrus.Errorf("Could not upload the file to the bucket. err: %v", err)
		return err
	}

	if err := wc.Close(); err != nil {
		logrus.Errorf("Could not close the write. err: %v", err)
		return err
	}

	return nil
}

func init() {
	flag.Parse()
}
