package main

import (
	"context"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type M3u8Collector struct {
	URL string
}

func (m *M3u8Collector) Collect(ctx context.Context) (*M3u8Downloader, error) {
	ctx, cancel := newChromedp(ctx, true)
	defer cancel()

	var (
		titleChan = make(chan string)
		m3u8Chan  = make(chan string)
		errChan   = make(chan error)
	)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			url := ev.Request.URL
			if strings.HasSuffix(url, ".m3u8") {
				m3u8Chan <- url
			}
		}
	})

	go func() {
		var title string
		if err := chromedp.Run(ctx, network.Enable(), chromedp.Navigate(m.URL), chromedp.Title(&title)); err != nil {
			errChan <- err
			return
		}
		titleChan <- title
	}()

	select {
	case url := <-m3u8Chan:
		return &M3u8Downloader{
			Title: <-titleChan,
			M3u8:  url,
		}, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
