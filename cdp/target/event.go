package target

import (
	"errors"

	"github.com/jacexh/guppeteer/cdp"
)

const (
	ReceivedMessageFromTarget = "Target.receivedMessageFromTarget"
	DetachedFromTarget        = "Target.detachedFromTarget"
)

type (
	EventReceivedMessageFromTarget struct {
		domain
	}
	ReceivedMessageFromTargetParams struct {
		SessionID SessionID `json:"sessionId"`
		Message   string    `json:"message"`
		TargetID  TargetID  `json:"targetId,omitempty"`
	}

	EventDetachedFromTarget struct {
		domain
	}

	DetachedFromTargetParams struct {
		SessionID SessionID `json:"sessionId"`
		TargetID  TargetID  `json:"targetId,omitempty"`
	}
)

func (rmft *EventReceivedMessageFromTarget) Name() string                       { return ReceivedMessageFromTarget }
func (rmft *EventReceivedMessageFromTarget) Load(d []byte) (interface{}, error) { return pHub(rmft, d) }

func (dftp *EventDetachedFromTarget) Name() string                       { return DetachedFromTarget }
func (dftp *EventDetachedFromTarget) Load(d []byte) (interface{}, error) { return pHub(dftp, d) }

func pHub(e cdp.Element, data []byte) (interface{}, error) {
	var params interface{}
	var err error
	switch e.Name() {
	case ReceivedMessageFromTarget:
		params = new(ReceivedMessageFromTargetParams)
	case DetachedFromTarget:
		params = new(DetachedFromTargetParams)
	default:
		return nil, errors.New("unsupported event: " + e.Name())
	}
	err = json.Unmarshal(data, params)
	return params, err

}
