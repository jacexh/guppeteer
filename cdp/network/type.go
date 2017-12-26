package network

const (
	domainName = "Network"
)

type (
	domain struct{}

	// LoaderID Unique loader identifier.
	LoaderID string
	// ConnectionType The underlying connection technology that the browser is supposedly using.
	// Allowed: none, cellular2g, cellular3g, cellular4g, bluetooth, ethernet, wifi, wimax, other.
	ConnectionType string
)

func (d domain) Domain() string { return domainName }
