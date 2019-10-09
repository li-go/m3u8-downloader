package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s {url}", os.Args[0])
		os.Exit(2)
	}

	collector := &M3u8Collector{URL: os.Args[1]}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	downloader, err := collector.Collect(ctx)
	if err != nil {
		log.Printf("fail to collect m3u8 url: %v", err)
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	if err := downloader.Download(ctx); err != nil {
		log.Printf("fail to download m3u8 url: %v", err)
	}
}
