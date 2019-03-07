// sample storage-quickstart creates a Google Cloud Storage bucket.

package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {

	files, err := list_all("bucket-messagebird-sample1", "voice-asterisk")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %s\n", files)

	// test */10
	files, err = list_all("bucket-messagebird-sample1", "voice-asterisk*/10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %s\n", files)

	// test */10
	files, err = list_all("bucket-messagebird-sample1", "voice-asterisk-(0[1-9]|[1-9][0-9])/10/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %s\n", files)

	// test */10
	files, err = list_all2("bucket-messagebird-sample1", "voice-asterisk[**]/10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %s\n", files)

	// files, err = list_regex("bucket-messagebird-sample1", "voice-asterisk", "asterisk-directory/voice-asterisk-(0[1-9]|[1-9][0-9])/10/*")
	files, err = list_regex("bucket-messagebird-sample1", "voice-asterisk", "asterisk-directory/voice-asterisk-*/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result list_regex: %s\n", files)

	files, err = list_date("bucket-messagebird-sample1", "asterisk-directory/voice-asterisk", time.Unix(1546171200, 0))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result list_date: %s\n", files)

}

func list_all2(bucketName string, prefix string) ([]string, error) {
	ctx := context.Background()

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client. err: %v", err)
		return nil, err
	}

	var files []string
	it := client.Bucket(bucketName).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		files = append(files, attrs.Name)
	}

	return files, nil
}

func list_date(bucketName string, prefix string, from time.Time) ([]string, error) {

	// get all files
	files, err := list_all(bucketName, prefix)
	if err != nil {
		return nil, err
	}
	fmt.Printf("tmp res: %s\n", files)

	// match
	// matchString := fmt.Sprintf("%s/cdr/%02d/%02d/%02d/*", "*/*", from.Year(), from.Month(), from.Day())
	// matchString := fmt.Sprintf("asterisk-directory/voice-asterisk-production-europe-west1-d-9/cdr/2018/9/30/cdr-20180930080001.csv")
	matchString := fmt.Sprintf("(.*)/cdr/%02d/%02d/%02d/(.*)", from.Year(), from.Month(), from.Day())
	fmt.Println(matchString)
	var res []string
	for _, file := range files {
		match, _ := regexp.MatchString(matchString, file)
		if match == true {
			res = append(res, file)
		}
	}

	return res, nil
}

func list_regex(bucketName string, prefix string, match string) ([]string, error) {
	files, err := list_all(bucketName, prefix)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, file := range files {
		fmt.Println(file)
		match, _ := regexp.MatchString(match, file)
		if match == true {
			fmt.Println("Matched!")
			res = append(res, file)
		}
	}

	return res, nil
}

func list_all(bucketName string, prefix string) ([]string, error) {
	ctx := context.Background()

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client. err: %v", err)
		return nil, err
	}

	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{
		Prefix: prefix,
	})

	var files []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, attrs.Name)
	}

	return files, nil

}
