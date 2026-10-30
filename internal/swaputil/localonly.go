package swaputil

import (
	"net/http"
	"strings"
)

// HeaderLocalOnly is the request header that restricts routing to local
// targets: models this instance serves itself, and peers with trust: local.
const HeaderLocalOnly = "X-Local-Only"

// LocalOnlyRejectMessage is the proxy log line and the JSON error message
// returned when a flagged request has no local target. The proxy must not
// fall through to a remote peer.
const LocalOnlyRejectMessage = "rejected: local_only, no local target available"

// LocalOnlyRequested reports whether the client asked to stay on local targets.
//
// The header's presence is enough, including an empty value. "0" and "false"
// (any case, surrounding space ignored) mean the flag is not set. Any other
// value counts as set. When the header is repeated, one non-false value is
// enough: a client that sent both "false" and "1" is asking for the gate.
func LocalOnlyRequested(h http.Header) bool {
	vals := h.Values(HeaderLocalOnly)
	if len(vals) == 0 {
		return false
	}
	for _, v := range vals {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "0", "false":
		default:
			return true
		}
	}
	return false
}
