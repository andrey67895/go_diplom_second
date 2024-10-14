package helpers

import (
	"bytes"
	"compress/gzip"
	"io"
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

func DeCompress(data []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	all, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return all, nil
}
