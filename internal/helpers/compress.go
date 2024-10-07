package helpers

import (
	"bytes"
	"compress/gzip"
)

func Compress(data []byte) []byte {

	var b bytes.Buffer

	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	return b.Bytes()
}
