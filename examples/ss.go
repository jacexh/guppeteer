package main

import (
	"fmt"
	"time"

	"github.com/jacexh/guppeteer"
)

func main() {
	cnx, err := guppeteer.NewConnection("ws://127.0.0.1:51909/devtools/browser/cdf1ab25-6a18-4f32-b9e1-e56313d6e7f6")
	if err != nil {
		panic(err)
	}
	browser := guppeteer.NewBrowser(cnx)
	page, err := browser.NewPage()
	start := time.Now()
	page.Goto("http://www.baidu.com", "")
	fmt.Println(time.Since(start).String())
}
