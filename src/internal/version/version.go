// Package version carries the single source of truth for the service
// identity + version reported by /health, /version and build_info metrics.
// Keep Version in sync with service.yaml and the git tag.
package version

const (
	// ServiceID is the fleet registry id for this service.
	ServiceID = "country-iso-matcher"

	// Version is the released service version (no leading "v").
	Version = "0.1.0"
)
