package update

import (
	"runtime"

	"github.com/blang/semver"
)

// A Release describes a specific release of the software.
type Release struct {
	ID       string
	Changelog string
	Version  semver.Version
	Variants []Variant
}

// GetVariant fetches the variant entry from this release
// which matches the provided variant definition.
func (r *Release) GetVariant(variant *Variant) *Variant {
	for _, v := range r.Variants {
		if v.Equals(variant) {
			return &v
		}
	}

	return nil
}

// Latest returns the latest from a collection of releases.
func Latest(releases []Release) *Release {
	if releases == nil || len(releases) == 0 {
		return nil
	}

	latest := &releases[0]
	for _, release := range releases {
		if release.Version.GT(latest.Version) {
			latest = &release
		}
	}

	return latest
}

// LatestUpdate gets the latest release which is also an update
// over the currently running application version.
func LatestUpdate(releases []Release, appVersion string) *Release {
	latest := Latest(releases)

	if latest == nil {
		return nil
	}

	v, err := semver.Parse(appVersion)

	// Treat errors as an updatable release (so that we don't have
	// broken or dev clients in the wild unable to update).
	if err != nil || latest.Version.GT(v) {
		return latest
	}

	return nil
}

// A Variant describes a platform specific build of the software.
type Variant struct {
	ID       string
	Platform string
	Arch     string
}

// Equals determines whether two variants are equal to one another or not.
func (v *Variant) Equals(other *Variant) bool {
	return v.Arch == other.Arch && v.Platform == other.Platform
}

// MyPlatform returns the variant representing your current platform.
func MyPlatform() *Variant {
	return &Variant{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}
}
