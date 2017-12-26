package cdp

import (
	"fmt"
	"sync/atomic"

	"github.com/json-iterator/go"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	globalCounter int64
)

type (
	// Element method/event 不具备完整method、event能力
	Element interface {
		Domain() string
		Name() string
	}

	// Event CDP定义event
	Event interface {
		Load([]byte) (interface{}, error) // 反序列化params字段，返回对应的params对象
		Element
	}

	// Method CDP定义的method
	Method interface {
		Dump() ([]byte, error) // Dump: 序列化后作为Message.params字段 Load: 反序列化result，返回对应的returns对象
		Event
	}

	// Message CDP消息结构体
	Message struct {
		ID     int64               `json:"id,omitempty"`
		Method string              `json:"method,omitempty"`
		Params jsoniter.RawMessage `json:"params,omitempty"`
		Result jsoniter.RawMessage `json:"result,omitempty"`
		Err    *MessageError       `json:"error,omitempty"`
	}

	// MessageError Message的Error字段
	MessageError struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
)

// HasError 是否存在错误
func (msg *Message) HasError() bool {
	if msg.Err != nil {
		return true
	}
	return false
}

// GetError 获取error
func (msg *Message) GetError() error {
	return msg.Err
}

func (e *MessageError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Message)
}

func NewMessage(m Method) (*Message, error) {
	data, err := m.Dump()
	if err != nil {
		return nil, err
	}
	return &Message{ID: atomic.AddInt64(&globalCounter, 1), Method: m.Name(), Params: data}, nil
}

// NewAndDumpMessage 创建并序列化Message对象
func NewAndDumpMessage(m Method) (*Message, []byte, error) {
	msg, err := NewMessage(m)
	if err != nil {
		return msg, nil, err
	}
	data, err := json.Marshal(msg)
	return msg, data, err
}
