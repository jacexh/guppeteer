package guppeteer

import (
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
	_, err := s.Execute(&page.MethodEnable{})
	if err != nil {
		return nil, err
	}

	_, err = s.Execute(&network.MethodEnable{})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Page) Goto(url, referrer string) (*page.Frame, error) {
	nav := &page.MethodNavigate{URL: url, Referrer: referrer}

	retFrame := new(page.Frame)
	var retErr error
	sub := &Subscriber{}
	sub.Subscribe("Page.frameStoppedLoading", nil)
	sub.Subscribe("Page.frameNavigated", func(d []byte) {
		ev := &page.EventFrameNavigated{}
		ret, err := ev.Load(d)
		if err != nil {
			retErr = err
		} else {
			retFrame = ret.(*page.FrameNavigatedParams).Frame
		}
	})
	defaultEventloop.Register(p.session.ID, sub)
	defer defaultEventloop.Cancel(p.session.ID)

	_, err := p.session.Execute(nav)
	if err != nil {
		return nil, err
	}
	sub.WaitUtilPublished()

	if retErr != nil {
		return nil, retErr
	}
	return retFrame, nil
}
