package page

import (
	"errors"

	"github.com/jacexh/guppeteer/cdp"
)

const (
	FrameNavigated      = "Page.frameNavigated"      // Fired once navigation of the frame has completed. Frame is now associated with the new loader.
	FrameStoppedLoading = "Page.frameStoppedLoading" // Fired when frame has stopped loading.
)

type (
	EventFrameNavigated struct {
		domain
	}
	FrameNavigatedParams struct {
		Frame *Frame `json:"frame"`
	}

	EventFrameStoppedLoading struct {
		domain
	}
	FrameStoppedLoadingParams struct {
		FrameID FrameID `json:"frameId"`
	}
)

func (fn *EventFrameNavigated) Name() string                       { return FrameNavigated }
func (fn *EventFrameNavigated) Load(d []byte) (interface{}, error) { return pHub(fn, d) }

func (fsl *EventFrameStoppedLoading) Name() string                       { return FrameStoppedLoading }
func (fsl *EventFrameStoppedLoading) Load(d []byte) (interface{}, error) { return pHub(fsl, d) }

func pHub(e cdp.Element, data []byte) (interface{}, error) {
	var p interface{}
	var err error
	switch e.Name() {
	case FrameNavigated:
		p = new(FrameNavigatedParams)
	case FrameStoppedLoading:
		p = new(FrameStoppedLoadingParams)
	default:
		return nil, errors.New("unsupported event: " + e.Name())
	}
	err = json.Unmarshal(data, p)
	return p, err
}
