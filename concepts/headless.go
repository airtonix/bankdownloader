package main

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte

	if err := chromedp.Run(ctx, debug(&buf)); err != nil {
		fmt.Println(err)
	} else if len(buf) == 0 {
		fmt.Println("no data from print")
	}
}

func debug(res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(1280, 768),
		chromedp.Navigate("https://google.ca"),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}
