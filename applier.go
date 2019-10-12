package update

import (
	"context"
	"io"
	"os"
)

// The Applier is responsible for applying each stage of the update operation.
// It does not pass state between phases or attempt to manage the process.
type Applier struct {
}

// Prepare executes the preparation phase of an upgrade by downloading a release
// to the specified location.
func (a *Applier) Prepare(ctx context.Context, source Source, release *Release, variant *Variant, dest string) error {
	err := deleteFile(ctx, dest)
	if err != nil {
		return err
	}

	if variant == nil {
		variant = MyPlatform()
	}

	data, err := source.Download(release, variant)
	if err != nil {
		return err
	}
	defer data.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(f, data)
	if err != nil {
		return err
	}

	return nil
}

// Replace executes the replacement phase of an upgrade by copying the source
// file over the destination file.
func (a *Applier) Replace(ctx context.Context, src, dest string) error {
	err := deleteFile(ctx, dest)
	if err != nil {
		return err
	}

	err = copyFile(src, dest)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup removes the files which were used to execute the update operation.
func (a *Applier) Cleanup(ctx context.Context, src string) error {
	return deleteFile(ctx, src)
}
