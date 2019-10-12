package update

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubSource(t *testing.T) {
	s := NewGitHubSource("sierrasoftworks/git-tool", "v", "git-tool-")
	require.NotNil(t, s, "it should return a non-nil source")

	t.Run("Releases()", func(t *testing.T) {
		rs, err := s.Releases()
		require.NoError(t, err, "it should not return an error")

		assert.NotEmpty(t, rs, "it should return the list of releases")

		totalVariants := 0

		for _, r := range rs {
			assert.NotNil(t, r, "the release entries should not be nil")
			assert.Contains(t, r.ID, r.Version.String(), "the version should be derived from the tag")

			totalVariants += len(r.Variants)

			if len(r.Variants) != 0 {
				plats := map[string]struct{}{}
				arches := map[string]struct{}{}

				for _, v := range r.Variants {
					assert.NotNil(t, v, "the variant should not be nil")
					assert.NotEmpty(t, v.ID, "the variant ID should not be nil")
					assert.NotEmpty(t, v.Platform, "the platform should not be empty")
					assert.NotEmpty(t, v.Arch, "the arch should not be empty")

					plats[v.Platform] = struct{}{}
					arches[v.Arch] = struct{}{}
				}

				assert.Greater(t, len(plats), 1, "there should be at least one platform supported")
				assert.Greater(t, len(arches), 1, "there should be at least one architecture supported")
			}
		}

		assert.Greater(t, totalVariants, 0, "we should have tracked at least one variant")
	})

	t.Run("Download()", func(t *testing.T) {
		rs, err := s.Releases()
		require.NoError(t, err, "it should not return an error when fetching the release list")

		release := Latest(rs)
		require.NotNil(t, release, "we should be able to find the latest release")

		data, err := s.Download(release, MyPlatform())
		assert.NoError(t, err, "we should not receive an error downloading the release")
		assert.NotNil(t, data, "the data stream should not be nil")

		defer data.Close()

		buf := make([]byte, 512)
		length := 0

		for {
			read, err := data.Read(buf)
			length += read
			if err == io.EOF {
				break
			}
		}

		assert.Greater(t, length, 0, "we should have downloaded the release variant")
	})
}
