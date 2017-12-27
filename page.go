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
	p.session.invoke(&page.MethodEnable{}, cb1)
	select {
	case err := <-cb1.WaitError():
		return nil, err
	case <-cb1.WaitResult():
	}

	cb2 := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb2)
	p.session.invoke(&network.MethodEnable{}, cb2)
	select {
	case err := <-cb2.WaitError():
		return nil, err
	case <-cb2.WaitResult():
	}
	return p, nil
}

func (p *Page) Goto(url, referrer string) error {
	nav := &page.MethodNavigate{URL: url, Referrer: referrer}
	cb := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb)
	sub := &Subscriber{}
	sub.Subscribe("Page.frameStoppedLoading", func(d []byte) interface{} { return nil })
	defaultEventloop.Register(p.session.ID, sub)
	defer defaultEventloop.Cancel(p.session.ID)

	ret, err := p.session.Execute(nav)
	if err != nil {
		return err
	}
	sub.WaitUtilPublished()
	fmt.Println(ret.(*page.NavigateReturns))
	return nil
}
