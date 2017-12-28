package guppeteer

import "github.com/jacexh/guppeteer/cdp/target"

type (
	Browser struct {
		cnx *Connection
	}
)

func NewBrowser(cnx *Connection) *Browser {
	return &Browser{cnx: cnx}
}

func (b *Browser) NewPage() (*Page, error) {
	method := &target.MethodCreateTarget{URL: "about:blank"}
	r, err := b.cnx.Execute(method)
	if err != nil {
		return nil, err
	}
	ret := r.(*target.CreateTargetReturns)
	ss, err := b.cnx.CreateSession(ret.TargetID)
	if err != nil {
		return nil, err
	}
	return CreatePage(ss)
}
