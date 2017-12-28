package main

import (
	"fmt"
	"time"

	"github.com/jacexh/guppeteer"
)

const (
	wsAddr = "ws://127.0.0.1:32769/devtools/browser/bceaadd9-4af3-4845-8f7d-56aa8a25eee5"
)

func main() {
	cnx, err := guppeteer.NewConnection(wsAddr)
	if err != nil {
		panic(err)
	}
	browser := guppeteer.NewBrowser(cnx)
	page, err := browser.NewPage()
	start := time.Now()
	f, err := page.Goto("http://www.baidu.com", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start).String())
	fmt.Println(f.URL)
}
