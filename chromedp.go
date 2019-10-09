package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func newChromedp(ctx context.Context, headless bool) (context.Context, context.CancelFunc) {
	var opts []chromedp.ExecAllocatorOption
	for _, opt := range chromedp.DefaultExecAllocatorOptions {
		opts = append(opts, opt)
	}
	if !headless {
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
		)
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return ctx, func() {
		cancel()
		allocCancel()
	}
}
