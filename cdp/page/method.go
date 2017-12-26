package page

import (
	"errors"

	"github.com/jacexh/guppeteer/cdp"
	"github.com/jacexh/guppeteer/cdp/network"
	"github.com/json-iterator/go"
)

const (
	CaptureScreenshot = "Page.captureScreenshot" // Capture page screenshot
	Enable            = "Page.enable"            // Enables page domain notifications
	Disable           = "Page.disable"           // Disables page domain notifications.
	Navigate          = "Page.navigate"          // Navigates current page to the given URL.
	Reload            = "Page.reload"            // Reloads given page optionally ignoring the cache.
	StopLoading       = "Page.stopLoading"       // Force the page stop all navigations and pending resource fetches.
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type (
	MethodCaptureScreenshot struct {
		domain
		Format      string    `json:"format,omitempty"`
		Quality     int       `json:"quality,omitempty"`
		Clip        *Viewport `json:"clip,omitempty"`
		FromSurface bool      `json:"fromSurface,omitempty"`
	}
	CaptureScreenshotReturns struct {
		Data string `json:"data"`
	}

	MethodEnable struct {
		domain
	}
	EnableReturns struct{}

	MethodDisable struct {
		domain
	}
	DisableReturn struct{}

	MethodNavigate struct {
		domain
		URL            string         `json:"url"`
		Referrer       string         `json:"referrer,omitempty"`
		TransitionType TransitionType `json:"transitionType,omitempty"`
	}
	NavigateReturns struct {
		FrameID   FrameID          `json:"frameId"`
		LoaderID  network.LoaderID `json:"loaderId,omitempty"`
		ErrorText string           `json:"errorText,omitempty"`
	}

	MethodReload struct {
		domain
		IgnoreCache            bool   `json:"ignoreCache,omitempty"`
		ScriptToEvaluateOnLoad string `json:"scriptToEvaluateOnLoad,omitempty"`
	}
	ReloadReturns struct{}

	MethodStopLoading struct {
		domain
	}
	StopLoadingReturns struct{}
)

func (cs *MethodCaptureScreenshot) Name() string                       { return CaptureScreenshot }
func (cs *MethodCaptureScreenshot) Dump() ([]byte, error)              { return json.Marshal(cs) }
func (cs *MethodCaptureScreenshot) Load(d []byte) (interface{}, error) { return retHub(cs, d) }

func (en *MethodEnable) Name() string                       { return Enable }
func (en *MethodEnable) Dump() ([]byte, error)              { return json.Marshal(en) }
func (en *MethodEnable) Load(d []byte) (interface{}, error) { return retHub(en, d) }

func (ds *MethodDisable) Name() string                       { return Disable }
func (ds *MethodDisable) Dump() ([]byte, error)              { return json.Marshal(ds) }
func (ds *MethodDisable) Load(d []byte) (interface{}, error) { return retHub(ds, d) }

func (n *MethodNavigate) Name() string                       { return Navigate }
func (n *MethodNavigate) Dump() ([]byte, error)              { return json.Marshal(n) }
func (n *MethodNavigate) Load(d []byte) (interface{}, error) { return retHub(n, d) }

func (r *MethodReload) Name() string                       { return Reload }
func (r *MethodReload) Dump() ([]byte, error)              { return json.Marshal(r) }
func (r *MethodReload) Load(d []byte) (interface{}, error) { return retHub(r, d) }

func (sl *MethodStopLoading) Name() string                       { return StopLoading }
func (sl *MethodStopLoading) Dump() ([]byte, error)              { return json.Marshal(sl) }
func (sl *MethodStopLoading) Load(d []byte) (interface{}, error) { return retHub(sl, d) }

func retHub(e cdp.Element, d []byte) (interface{}, error) {
	var ret interface{}
	var err error
	switch e.Name() {
	case CaptureScreenshot:
		ret = new(CaptureScreenshotReturns)
	case Enable:
		return EnableReturns{}, nil
	case Disable:
		return DisableReturn{}, nil
	case Reload:
		return ReloadReturns{}, nil
	case StopLoading:
		return StopLoadingReturns{}, nil
	case Navigate:
		ret = new(NavigateReturns)
	default:
		return nil, errors.New("unsupported method: " + e.Name())
	}

	err = json.Unmarshal(d, ret)
	return ret, err
}
