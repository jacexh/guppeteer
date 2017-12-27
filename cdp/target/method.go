package target

import (
	"errors"

	"github.com/jacexh/guppeteer/cdp"
	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type (
	// MethodActivateTarget Activates (focuses) the target
	MethodActivateTarget struct {
		domain
		TargetID TargetID `json:"targetId"`
	}
	// ActivateTargetReturns result for page.activateTarget
	ActivateTargetReturns struct{}

	MethodAttachToTarget struct {
		domain
		TargetID TargetID `json:"targetId"`
	}
	AttachToTargetReturns struct {
		SessionID SessionID `json:"sessionId"`
	}

	MethodCreateTarget struct {
		domain
		URL                     string           `json:"url"`
		Width                   int              `json:"width,omitempty"`
		Height                  int              `json:"height,omitempty"`
		BrowserContextID        BrowserContextID `json:"browserContextId,omitempty"`
		EnableBeginFrameControl bool             `json:"enableBeginFrameControl,omitempty"`
	}
	CreateTargetReturns struct {
		TargetID TargetID `json:"targetId"`
	}

	MethodSendMessageToTarget struct {
		domain
		Message   string    `json:"message"`
		SessionID SessionID `json:"sessionId,omitempty"`
		TargetID  TargetID  `json:"targetId,omitempty"`
	}
	SendMessageToTargetReturns struct{}
)

const (
	ActivateTarget      = "Target.activateTarget"
	AttachToTarget      = "Target.attachToTarget"
	CreateTarget        = "Target.createTarget"
	SendMessageToTarget = "Target.sendMessageToTarget"
)

func (at *MethodActivateTarget) Name() string                          { return ActivateTarget }
func (at *MethodActivateTarget) Load(data []byte) (interface{}, error) { return retHub(at, data) }
func (at *MethodActivateTarget) Dump() ([]byte, error)                 { return json.Marshal(at) }

func (att *MethodAttachToTarget) Name() string                          { return AttachToTarget }
func (att *MethodAttachToTarget) Load(data []byte) (interface{}, error) { return retHub(att, data) }
func (att *MethodAttachToTarget) Dump() ([]byte, error)                 { return json.Marshal(att) }

func (ct *MethodCreateTarget) Name() string                          { return CreateTarget }
func (ct *MethodCreateTarget) Load(data []byte) (interface{}, error) { return retHub(ct, data) }
func (ct *MethodCreateTarget) Dump() ([]byte, error)                 { return json.Marshal(ct) }

func (smtt *MethodSendMessageToTarget) Name() string                       { return SendMessageToTarget }
func (smtt *MethodSendMessageToTarget) Dump() ([]byte, error)              { return json.Marshal(smtt) }
func (smtt *MethodSendMessageToTarget) Load(d []byte) (interface{}, error) { return retHub(smtt, d) }

func retHub(em cdp.Element, data []byte) (interface{}, error) {
	var ret interface{}
	var err error

	switch em.Name() {
	case ActivateTarget:
		return ActivateTargetReturns{}, nil
	case SendMessageToTarget:
		return SendMessageToTargetReturns{}, nil
	case AttachToTarget:
		ret = new(AttachToTargetReturns)
	case CreateTarget:
		ret = new(CreateTargetReturns)
	default:
		return nil, errors.New("unsupported method: " + em.Name())
	}
	err = json.Unmarshal(data, ret)
	return ret, err
}
