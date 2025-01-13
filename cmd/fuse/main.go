package main

import (
	"context"
	"fmt"
	"os"

	"github.com/billziss-gh/cgofuse/fuse"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
)

// Created bucket by the following command:
// ```shell
// gcloud storage buckets create gs://learn_systems_programming_with_go_20250113 --location=asia-northeast1
// ```
// Uploaded some files by the following command:
// ```shell
// gcloud storage cp test.1.txt test.2.txt gs://learn_systems_programming_with_go_20250113
// ```
// const (
// 	bucketURL = "gs://learn_systems_programming_with_go_20250113"
// )

type cloudFS struct {
	fuse.FileSystemBase
	bucket *blob.Bucket
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s [bucket-URL] [mount-point]\n", os.Args[0])
		os.Exit(1)
	}
	bucketURL := os.Args[1]

	ctx := context.Background()
	b, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		panic(err)
	}
	defer b.Close()
	obj, err := b.List(&blob.ListOptions{Prefix: ""}).Next(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("obj: %#v\n", obj)

	// Not works on my machine due to privacy settings
}
