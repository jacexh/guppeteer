package guppeteer

import (
	"errors"
	"sync"
	"sync/atomic"
)

type (
	Subscriber struct {
		receivers sync.Map
		wg        sync.WaitGroup
	}

	Receiver struct {
		received int32
		event    string
		f        func([]byte) interface{}
		ret      interface{}
	}
)

func (sub *Subscriber) Subscribe(event string, f func([]byte) interface{}) error {
	if _, loaded := sub.receivers.LoadOrStore(event, NewReceiver(event, f)); !loaded {
		sub.wg.Add(1)
		return nil
	}
	return errors.New("duplicated event")
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

func NewReceiver(event string, f func([]byte) interface{}) *Receiver {
	return &Receiver{event: event, f: f}
}

func (rc *Receiver) Receive(d []byte) (interface{}, error) {
	if atomic.CompareAndSwapInt32(&rc.received, 0, 1) {
		rc.ret = rc.f(d)
		return rc.ret, nil
	}
	return nil, errors.New("received message yet")
}

func (rc *Receiver) Returns() interface{} {
	return rc.ret
}
