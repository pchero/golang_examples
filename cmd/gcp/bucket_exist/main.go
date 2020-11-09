package main

// go run main.go -bucket_name <bucket name> -file <filename>

import (
	"context"
	"flag"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

var bucketName = flag.String("bucket_name", "sample_bucket", "Bucket name for gcp storage")
var filename = flag.String("file", "test.go", "Filename for upload")

func main() {
	ctx := context.Background()

	flag.Parse()

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Errorf("Failed to create client. err: %v", err)
		return
	}

	f := client.Bucket(*bucketName).Object(*filename)
	attrs, err := f.Attrs(ctx)
	if err != nil {
		logrus.Errorf("Could not get fie attr. err: %v", err)
		return
	}

	logrus.Infof("Bucket: %v\n", attrs.Bucket)
	logrus.Infof("CacheControl: %v\n", attrs.CacheControl)
	logrus.Infof("ContentDisposition: %v\n", attrs.ContentDisposition)
	logrus.Infof("ContentEncoding: %v\n", attrs.ContentEncoding)
	logrus.Infof("ContentLanguage: %v\n", attrs.ContentLanguage)
	logrus.Infof("ContentType: %v\n", attrs.ContentType)
	logrus.Infof("Crc32c: %v\n", attrs.CRC32C)
	logrus.Infof("Generation: %v\n", attrs.Generation)
	logrus.Infof("KmsKeyName: %v\n", attrs.KMSKeyName)
	logrus.Infof("Md5Hash: %v\n", attrs.MD5)
	logrus.Infof("MediaLink: %v\n", attrs.MediaLink)
	logrus.Infof("Metageneration: %v\n", attrs.Metageneration)
	logrus.Infof("Name: %v\n", attrs.Name)
	logrus.Infof("Size: %v\n", attrs.Size)
	logrus.Infof("StorageClass: %v\n", attrs.StorageClass)
	logrus.Infof("TimeCreated: %v\n", attrs.Created)
	logrus.Infof("Updated: %v\n", attrs.Updated)
	logrus.Infof("Event-based hold enabled? %t\n", attrs.EventBasedHold)
	logrus.Infof("Temporary hold enabled? %t\n", attrs.TemporaryHold)
	logrus.Infof("Retention expiration time %v\n", attrs.RetentionExpirationTime)
	logrus.Infof("Custom time %v\n", attrs.CustomTime)
	logrus.Infof("\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		logrus.Infof("\t%v = %v\n", key, value)
	}
}
