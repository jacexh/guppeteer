package network

import (
	"errors"

	"github.com/jacexh/guppeteer/cdp"
	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	// Enable network tracking, network events will now be delivered to the client
	Enable = "Network.enable" // Enables .
	// Disable network tracking, prevents network events from being sent to the client
	Disable = "Network.disable"
	// EmulateNetworkConditions Activates emulation of network conditions.
	EmulateNetworkConditions = "Nework.emulateNetworkConditions"
)

type (
	MethodEnable struct {
		domain
		MaxTotalBufferSize    int `json:"maxTotalBufferSize,omitempty"`
		MaxResourceBufferSize int `json:"maxResourceBufferSize,omitempty"`
	}
	EnableReturns struct{}

	MethodDisable struct {
		domain
	}
	DisableReturns struct{}

	MethodEmulateNetworkConditions struct {
		domain
		Offline            bool           `json:"offline"`
		Latency            int            `json:"latency"`            // ms
		DownloadThroughput int            `json:"downloadThroughput"` // Maximal aggregated download throughput (bytes/sec). -1 disables download throttling.
		UploadThroughput   int            `json:"uploadThroughput"`   // Maximal aggregated upload throughput (bytes/sec). -1 disables upload throttling
		ConnectionType     ConnectionType `json:"connectionType,omitempty"`
	}
	EmulateNetworkConditionsReturns struct{}
)

func (e *MethodEnable) Name() string                       { return Enable }
func (e *MethodEnable) Dump() ([]byte, error)              { return json.Marshal(e) }
func (e *MethodEnable) Load(d []byte) (interface{}, error) { return retHub(e, d) }

func (db *MethodDisable) Name() string                       { return Disable }
func (db *MethodDisable) Dump() ([]byte, error)              { return json.Marshal(db) }
func (db *MethodDisable) Load(d []byte) (interface{}, error) { return retHub(db, d) }

func (enc *MethodEmulateNetworkConditions) Name() string                       { return EmulateNetworkConditions }
func (enc *MethodEmulateNetworkConditions) Dump() ([]byte, error)              { return json.Marshal(enc) }
func (enc *MethodEmulateNetworkConditions) Load(d []byte) (interface{}, error) { return retHub(enc, d) }

func retHub(em cdp.Element, d []byte) (interface{}, error) {
	switch em.Name() {
	case Enable:
		return EnableReturns{}, nil
	case Disable:
		return DisableReturns{}, nil
	case EmulateNetworkConditions:
		return EmulateNetworkConditionsReturns{}, nil
	default:
		return nil, errors.New("unsupported method: " + em.Name())
	}
}
