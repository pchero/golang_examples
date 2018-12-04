// sample storage-quickstart creates a Google Cloud Storage bucket.

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "smiling-axiom-219521"
	//	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	it := client.Buckets(ctx, projectID)
	fmt.Println("Buckets:")
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(battrs.Name)
	}
}
