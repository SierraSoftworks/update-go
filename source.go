package update

import "io"

// A Source provides both information on the availability of versions
// as well as the ability to download a specific version on request.
type Source interface {
	Releases() ([]Release, error)
	Download(release *Release, variant *Variant) (io.ReadCloser, error)
}
