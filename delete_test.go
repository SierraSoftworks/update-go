package update

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteFile(t *testing.T) {
	tempDir := os.TempDir()

	t.Run("When the file doesn't exist", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "updatego.test.deleteFile.nonexistent")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := deleteFile(ctx, filePath)
		assert.NoError(t, err, "there should be no error removing the file")

		f, err := os.Open(filePath)
		assert.True(t, os.IsNotExist(err), "the file should not exist")
		defer f.Close()
	})

	t.Run("When the file exists", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "updatego.test.deleteFile.existent")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		f, err := os.OpenFile(filePath, os.O_CREATE, os.ModePerm)
		require.NoError(t, err)

		defer f.Close()
		defer os.Remove(filePath)

		fileUnlocked := false

		go func() {
			time.Sleep(400 * time.Millisecond)
			fileUnlocked = true
			f.Close()
		}()

		err = deleteFile(ctx, filePath)
		assert.NoError(t, err, "there should be no error removing the file")
		assert.True(t, fileUnlocked, "the file should only be removed once it is unlocked")

		f2, err := os.Open(filePath)
		assert.True(t, os.IsNotExist(err), "the file should not exist")
		defer f2.Close()
	})
}
