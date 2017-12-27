package guppeteer

import (
	"fmt"

	"github.com/jacexh/guppeteer/cdp/network"
	"github.com/jacexh/guppeteer/cdp/page"
)

type (
	Page struct {
		session *Session
	}
)

func CreatePage(s *Session) (*Page, error) {
	p := &Page{s}
	cb1 := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb1)
	p.session.CallMethod(&page.MethodEnable{}, cb1)
	select {
	case err := <-cb1.WaitError():
		return nil, err
	case <-cb1.WaitResult():
	}

	cb2 := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb2)
	p.session.CallMethod(&network.MethodEnable{}, cb2)
	select {
	case err := <-cb2.WaitError():
		return nil, err
	case <-cb2.WaitResult():
	}
	return p, nil
}

func (p *Page) Goto(url, referrer string) {
	nav := &page.MethodNavigate{URL: url, Referrer: referrer}
	cb := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb)
	sub := &Subscriber{}
	sub.Subscribe("Page.frameStoppedLoading", func(d []byte) interface{} {
		fmt.Println(string(d))
		fmt.Println("hello world")
		return nil
	})
	defaultEventloop.Register(p.session.ID, sub)

	err := p.session.CallMethod(nav, cb)
	if err != nil {
		cb.Deleted()
		return
	}

	select {
	case err = <-cb.WaitError():
		return
	case <-cb.WaitResult():
		// todo:
	}
	sub.WaitUtilPublished()
}
