package ptrss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/moviegeek/ptrss/pkg/parser"
	"github.com/moviegeek/ptrss/pkg/store"

	"cloud.google.com/go/storage"
	"gocloud.dev/blob"

	//use gcs
	_ "gocloud.dev/blob/gcsblob"
)

const (
	gcsBucketEnv  = "PTRSS_GCS_BUCKET"
	jsonObjectEnv = "PTRSS_JSON_FILENAME"
	xmlObjectEnv  = "PTRSS_XML_FILENAME"
)

//PubSubMessage the payload of the pub/sub event
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Config struct {
	BucketURL     string
	JsonObjectKey string
	XmlObjectKey  string
}

//UpdateRss is the entry point for gcloud function
func UpdateRss(ctx context.Context, m PubSubMessage) error {
	config, err := readEnv()
	if err != nil {
		return err
	}

	b, err := blob.OpenBucket(ctx, config.BucketURL)
	if err != nil {
		return fmt.Errorf("Failed to setup bucket: %s", err)
	}

	jsonFileWriter, err := b.NewWriter(ctx, config.JsonObjectKey, &blob.WriterOptions{
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}

	beforeWrite := func(as func(interface{}) bool) error {
		var sw *storage.Writer
		if as(&sw) {
			sw.PredefinedACL = "publicRead"
		}
		return nil
	}
	xmlFileWriter, err := b.NewWriter(ctx, config.XmlObjectKey, &blob.WriterOptions{
		ContentType: "application/rss+xml;charset=utf-8",
		BeforeWrite: beforeWrite,
	})
	if err != nil {
		return err
	}

	parser := parser.New()
	feedItems := []*gofeed.Item{}

	if hdcItems, err := parser.GetHDCFeedItems(); err != nil {
		fmt.Printf("failed to parse HDC rss feed, skip it")
	} else {
		log.Printf("got %d items from HDC", len(hdcItems))
		feedItems = append(feedItems, hdcItems...)
	}

	if ptItems, err := parser.GetPutaoFeedItems(); err != nil {
		fmt.Printf("failed to parse Putao rss feed, skip it")
	} else {
		log.Printf("got %d items from Putao", len(ptItems))
		feedItems = append(feedItems, ptItems...)
	}

	s := store.New()

	for _, item := range feedItems {
		s.AddFromFeedItem(item)
	}

	movies := s.Movies()

	err = json.NewEncoder(jsonFileWriter).Encode(movies)
	if err != nil {
		fmt.Printf("failed to store movies json file: %v\n", err)
	}
	fmt.Println("writing movies to json file")
	err = jsonFileWriter.Close()
	if err != nil {
		fmt.Printf("write failed: %v", err)
	}

	err = s.ToRss(xmlFileWriter)
	if err != nil {
		return fmt.Errorf("failed to store rss: %v", err)
	}
	fmt.Println("writing movies to rss xml file")
	err = xmlFileWriter.Close()
	if err != nil {
		fmt.Printf("failed to write xml file: %v", err)
		return err
	}

	return nil
}

func readEnv() (Config, error) {
	config := Config{}

	v := os.Getenv(gcsBucketEnv)
	if v == "" {
		return config, fmt.Errorf("gcs bucket is not set in environment variable %s", gcsBucketEnv)
	}
	config.BucketURL = v

	v = os.Getenv(jsonObjectEnv)
	if v == "" {
		return config, fmt.Errorf("gcs bucket is not set in environment variable %s", gcsBucketEnv)
	}
	config.JsonObjectKey = v

	v = os.Getenv(xmlObjectEnv)
	if v == "" {
		return config, fmt.Errorf("gcs bucket is not set in environment variable %s", gcsBucketEnv)
	}
	config.XmlObjectKey = v

	return config, nil
}
