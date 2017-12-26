package target

const (
	domainName = "Target"
)

type (
	// TargetID string
	TargetID string

	// SessionID Unique identifier of attached debugging session
	SessionID string

	// BrowserContextID .
	BrowserContextID string

	domain struct{}
)

func (d domain) Domain() string {
	return domainName
}
