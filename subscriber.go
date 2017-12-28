package guppeteer

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/jacexh/guppeteer/cdp/target"
)

type (
	Subscriber struct {
		receivers sync.Map
		wg        sync.WaitGroup
	}

	Receiver struct {
		received int32
		event    string
		f        func([]byte)
	}

	eventloop struct {
		sessions sync.Map
	}
)

var (
	defaultEventloop = &eventloop{}
)

func (el *eventloop) Register(sid target.SessionID, sub *Subscriber) {
	el.sessions.Store(sid, sub)
}

func (el *eventloop) Cancel(sid target.SessionID) {
	el.sessions.Delete(sid)
}

func (el *eventloop) Handle(sid target.SessionID, event string, d []byte) {
	if val, loaded := el.sessions.Load(sid); loaded {
		sub := val.(*Subscriber)
		sub.Handle(event, d)
	}
}

func (sub *Subscriber) Subscribe(event string, f func([]byte)) {
	sub.wg.Add(1)
	sub.receivers.Store(event, NewReceiver(event, f))
}

func (sub *Subscriber) Handle(event string, d []byte) {
	if val, ok := sub.receivers.Load(event); ok {
		go func(r *Receiver) {
			defer sub.wg.Done()
			r.Receive(d)
		}(val.(*Receiver))
	}
}

func (sub *Subscriber) WaitUtilPublished() map[string]*Receiver {
	sub.wg.Wait()
	ret := map[string]*Receiver{}
	sub.receivers.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value.(*Receiver)
		return true
	})
	return ret
}

func NewReceiver(event string, f func([]byte)) *Receiver {
	return &Receiver{event: event, f: f}
}

func (rc *Receiver) Receive(d []byte) error {
	if atomic.CompareAndSwapInt32(&rc.received, 0, 1) {
		if rc.f != nil {
			rc.f(d)
		}
		return nil
	}
	return errors.New("received message")
}
