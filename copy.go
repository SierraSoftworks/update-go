package update

import (
	"io"
	"os"
)

func copyFile(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	tgt, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}
	defer tgt.Close()

	err = tgt.Truncate(0)
	if err != nil {
		return err
	}

	_, err = io.Copy(tgt, src)
	if err != nil {
		return err
	}

	return nil
}
