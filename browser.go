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
	callback := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(callback)
	method := &target.MethodCreateTarget{URL: "about:blank"}
	err := b.cnx.invoke(method, callback)
	if err != nil {
		return nil, err
	}
	select {
	case err = <-callback.WaitError():
		return nil, err
	case data := <-callback.WaitResult():
		r, _ := method.Load(data)
		ret := r.(*target.CreateTargetReturns)
		ss, err := b.cnx.CreateSession(ret.TargetID)
		if err != nil {
			return nil, err
		}
		return CreatePage(ss)
	}
}
