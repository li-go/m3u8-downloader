package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type M3u8Downloader struct {
	Title string
	M3u8  string
}

func (d *M3u8Downloader) Mp4File() string {
	return d.Title + ".mp4"
}

func (d *M3u8Downloader) DownloadedSize() string {
	if info, err := os.Stat(d.Mp4File()); err == nil {
		return formatFileSize(info.Size())
	}
	return "0"
}

func (d *M3u8Downloader) Download(ctx context.Context) error {
	downloadedSize := d.DownloadedSize()
	if downloadedSize != "0" {
		log.Printf("%s already downloaded %s", d.Mp4File(), downloadedSize)
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("%s downloading %s", d.Mp4File(), d.DownloadedSize())
			case <-ctx.Done():
				break
			}
		}
	}()

	shCmd := fmt.Sprintf(`ffmpeg -i "%s" -vcodec copy -c copy -c:a aac "%s"`, d.M3u8, d.Mp4File())
	cmd := exec.CommandContext(ctx, "sh", "-c", shCmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Printf("%s downloaded %s", d.Mp4File(), d.DownloadedSize())
	return nil
}

func formatFileSize(size int64) string {
	units := []string{"G", "M", "K"}
	sizePerUnits := []int64{1_000_000_000, 1_000_000, 1_000}
	for i, unit := range units {
		if size >= sizePerUnits[i] {
			if size/sizePerUnits[i] >= 10 {
				return strconv.FormatInt(size/sizePerUnits[i], 10) + unit
			}
			return fmt.Sprintf("%.1f%s", float64(size)/float64(sizePerUnits[i]), unit)
		}
	}
	return strconv.FormatInt(size, 10)
}
