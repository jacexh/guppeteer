package main

import (
	"fmt"

	"time"

	"github.com/jacexh/guppeteer"
	"github.com/jacexh/guppeteer/cdp/page"
)

func main() {
	cnx, err := guppeteer.NewConnection("ws://127.0.0.1:51909/devtools/browser/cdf1ab25-6a18-4f32-b9e1-e56313d6e7f6")
	if err != nil {
		panic(err)
	}
	ss, err := cnx.CreateSession("(5727F98225E6591A24B9CEB6261DE79E)")
	if err != nil {
		panic(err)
	}

	callback := guppeteer.NewCallback()
	nav := &page.MethodNavigate{URL: "http://www.google.com"}
	ss.CallMethod(nav, callback)

	timeout := 5 * time.Second

	select {
	case ret := <-callback.WaitResult():
		r, _ := nav.Load(ret)
		//ret := r.(*page.NavigateReturns)
		fmt.Println(r.(*page.NavigateReturns).FrameID)
	case err := <-callback.WaitError():
		fmt.Println(err.Error())
	case <-time.NewTimer(timeout).C:
		callback.Deleted()
		fmt.Println("timeout")
	}
}
