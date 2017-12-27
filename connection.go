package guppeteer

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jacexh/guppeteer/cdp"
	"github.com/jacexh/guppeteer/cdp/target"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type (
	Connection struct {
		conn      *websocket.Conn
		callbacks sync.Map // 存放尚未返回结果的method： key -> message.ID
		sessions  sync.Map // 存放各个子session可能要用到的对象: key -> sessionID
	}

	Session struct {
		ID        target.SessionID
		TID       target.TargetID
		callbacks sync.Map
		parent    *Connection
	}
)

func NewConnection(url string) (*Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	c := &Connection{conn: conn}
	go c.receiveMessage()
	return c, nil
}

func (cnx *Connection) invoke(method cdp.Method, n *Callback) error {
	msg, data, err := cdp.NewAndDumpMessage(method)
	if err != nil {
		return err
	}
	return cnx.sendMessage(msg, data, n)
}

func (cnx *Connection) Execute(method cdp.Method) (interface{}, error) {
	cb := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb)
	err := cnx.invoke(method, cb)
	if err != nil {
		return nil, err
	}
	select {
	case err = <-cb.WaitError():
		return nil, err
	case data := <-cb.WaitResult():
		return method.Load(data)
	}
}

func (cnx *Connection) sendMessage(msg *cdp.Message, data []byte, n *Callback) error {
	var err error
	if msg == nil && data == nil {
		return errors.New("<nil>")
	} else if msg == nil && data != nil { // 此时直接写入消息
		return cnx.conn.WriteMessage(websocket.TextMessage, data)
	}

	if data == nil {
		data, err = json.Marshal(msg)
		if err != nil {
			return err
		}
	}

	err = cnx.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	Logger.Info("write", zap.ByteString("message", data))
	if msg.ID != 0 && n != nil { // 必须是成功写入消息后才存放
		cnx.callbacks.Store(msg.ID, n)
	}
	return nil
}

func (cnx *Connection) handleCallback(msg *cdp.Message) {
	val, loaded := cnx.callbacks.Load(msg.ID)
	if !loaded {
		return
	}
	defer cnx.callbacks.Delete(msg.ID)
	cb := val.(*Callback)
	if msg.HasError() {
		cb.SetError(msg.GetError())
	} else {
		cb.SetResult(msg.Result)
	}
}

func (cnx *Connection) handleEvent(msg *cdp.Message) {
	switch msg.Method {
	case target.ReceivedMessageFromTarget:
		rm := &target.EventReceivedMessageFromTarget{}
		params, _ := rm.Load(msg.Params)
		p := params.(*target.ReceivedMessageFromTargetParams)
		sid := p.SessionID
		sm := new(cdp.Message)
		err := json.Unmarshal([]byte(p.Message), sm)
		if err != nil {
			Logger.Error("unmarshal received message failed", zap.Error(err))
			return
		}
		if s, loaded := cnx.sessions.Load(sid); loaded {
			if sm.ID != 0 {
				s.(*Session).handleCallback(sm) // Session级别的method
			} else {
				defaultEventloop.Handle(sid, sm.Method, sm.Params)
			}
		}

	case target.DetachedFromTarget:
		d := &target.EventDetachedFromTarget{}
		p, _ := d.Load(msg.Params)
		params := p.(*target.DetachedFromTargetParams)
		if s, loaded := cnx.sessions.Load(params.SessionID); loaded {
			s.(*Session).Close()
			cnx.sessions.Delete(params.SessionID)
		}

	default:
	}
}

func (cnx *Connection) receiveMessage() error {
	for {
		mt, data, err := cnx.conn.ReadMessage()
		if err != nil {
			Logger.Error("occur error when reading message", zap.Error(err))
			return err
		}
		if mt != websocket.TextMessage {
			continue
		}
		msg := new(cdp.Message)
		err = json.Unmarshal(data, msg)
		if err != nil {
			Logger.Warn("unmarshal data failed", zap.Error(err))
			continue
		}
		Logger.Info("read", zap.ByteString("message", data))
		if msg.ID != 0 { // 处理method result
			cnx.handleCallback(msg)
		} else {
			cnx.handleEvent(msg)
		}
	}
}

func (cnx *Connection) CreateSession(tid target.TargetID) (*Session, error) {
	m := &target.MethodAttachToTarget{TargetID: tid}
	callback := defaultCallbackPool.Get()
	cnx.invoke(m, callback)
	defer func() { defaultCallbackPool.Put(callback) }()
	select {
	case err := <-callback.WaitError():
		return nil, err
	case data := <-callback.WaitResult():
		ret, err := m.Load(data)
		if err != nil {
			return nil, err
		}
		ss := &Session{ID: ret.(*target.AttachToTargetReturns).SessionID, TID: tid, parent: cnx}
		cnx.sessions.Store(ss.ID, ss)
		return ss, nil
	}
}

func (ss *Session) invoke(method cdp.Method, notifier *Callback) error {
	msg, data, err := cdp.NewAndDumpMessage(method)
	if err != nil {
		return err
	}
	err = ss.parent.invoke(&target.MethodSendMessageToTarget{Message: string(data), SessionID: ss.ID}, nil) // notifier 不需要传给Connection
	if err == nil && notifier != nil {
		ss.callbacks.Store(msg.ID, notifier)
	}
	return err
}

func (ss *Session) Execute(method cdp.Method) (interface{}, error) {
	cb := defaultCallbackPool.Get()
	defer defaultCallbackPool.Put(cb)
	err := ss.invoke(method, cb)
	if err != nil {
		return nil, err
	}

	select {
	case err = <-cb.WaitError():
		return nil, err
	case data := <-cb.WaitResult():
		return method.Load(data)
	}
}

func (ss *Session) handleCallback(msg *cdp.Message) {
	val, loaded := ss.callbacks.Load(msg.ID)
	if !loaded {
		return
	}
	defer ss.callbacks.Delete(msg.ID)
	cb := val.(*Callback)
	if msg.HasError() {
		cb.SetError(msg.GetError())
	} else {
		cb.SetResult(msg.Result)
	}
}

func (ss *Session) Close() error {
	// todo:
	return nil
}
