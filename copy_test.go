package update

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	tempDir := os.TempDir()

	sourcePath := filepath.Join(tempDir, "updatego.test.copyFile.source")
	destPath := filepath.Join(tempDir, "updatego.test.copyFile.dest")

	defer os.Remove(sourcePath)
	defer os.Remove(destPath)

	t.Run("When the source file doesn't exist", func(t *testing.T) {
		os.Remove(sourcePath)

		assert.Error(t, copyFile(sourcePath, destPath))
	})

	t.Run("When the source file exists", func(t *testing.T) {
		assert.NoError(t, ioutil.WriteFile(sourcePath, []byte("test"), os.ModePerm))
		os.Remove(destPath)

		assert.NoError(t, copyFile(sourcePath, destPath))

		data, err := ioutil.ReadFile(destPath)
		assert.NoError(t, err)
		assert.Equal(t, "test", string(data), "it should have written the right data to the target file")

		data, err = ioutil.ReadFile(sourcePath)
		assert.NoError(t, err)
		assert.Equal(t, "test", string(data), "it should not have modified the source file data")
	})

	t.Run("When the source and destination files exists", func(t *testing.T) {
		require.NoError(t, ioutil.WriteFile(sourcePath, []byte("test"), os.ModePerm))
		require.NoError(t, ioutil.WriteFile(destPath, []byte("not overwritten"), os.ModePerm))

		assert.Error(t, copyFile(sourcePath, destPath))

		data, err := ioutil.ReadFile(destPath)
		assert.NoError(t, err)
		assert.Equal(t, "not overwritten", string(data), "it should not have modified the target file data")

		data, err = ioutil.ReadFile(sourcePath)
		assert.NoError(t, err)
		assert.Equal(t, "test", string(data), "it should not have modified the source file data")
	})
}
