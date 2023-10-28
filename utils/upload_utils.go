package utils

import (
	"io"
	"os"
)

func SaveFileWithBuffer(fileObj io.Reader, out *os.File) error {
	buf := make([]byte, 1024)
	for {
		n, err := fileObj.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := out.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}
