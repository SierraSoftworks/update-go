package update

import (
	"context"
	"os"
	"time"

	"github.com/SierraSoftworks/rates"
)

func deleteFile(ctx context.Context, path string) error {
	bucket := rates.NewBucket(&rates.BucketConfig{
		MaxSize:  5,
		FillRate: 0.1,
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-bucket.TakeWhenAvailable():
			err := os.Remove(path)
			if err == nil || os.IsNotExist(err) {
				return nil
			}
		}

		// Don't poll more frequently than once per 250ms
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
}
