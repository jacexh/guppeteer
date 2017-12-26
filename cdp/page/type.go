package page

import "github.com/jacexh/guppeteer/cdp/network"

const (
	domainName = "Page"
)

type (
	domain struct{}

	// Viewport  for capturing screenshot.
	Viewport struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
		Scale  int `json:"scale"`
	}

	TransitionType string // allowed values: link, typed, auto_bookmark, auto_subframe, manual_subframe, generated, auto_toplevel, form_submit, reload, keyword, keyword_generated, other.

	FrameID string // Unique frame identifier
	Frame   struct {
		ID             string           `json:"id"`
		ParentID       string           `json:"parentId,omitempty"`
		LoaderID       network.LoaderID `json:"loaderId,omitempty"`
		Name           string           `json:"name,omitempty"`
		URL            string           `json:"url"`
		SecurityOrigin string           `json:"securityOrigin"`
		MimeType       string           `json:"mimeType"`
		UnreachableURL string           `json:"unreachableUrl,omitempty"`
	}
)

func (d domain) Domain() string {
	return domainName
}
