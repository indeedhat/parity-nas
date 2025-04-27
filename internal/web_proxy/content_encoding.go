package webproxy

import (
	"bytes"
	"compress/gzip"
	"io"
)

// decompress decompresses the given data stream with the provided encoding algo
func decompress(data io.Reader, algo string) (io.Reader, error) {
	switch algo {
	case "gzip":
		return gzip.NewReader(data)
	default:
		return data, nil
	}
}

// compress compresses the given data stream with the provided encoding algo
func compress(data io.Reader, algo string) (io.Reader, error) {
	switch algo {
	case "gzip":
		var buf bytes.Buffer

		writer := gzip.NewWriter(&buf)
		io.Copy(writer, data)
		writer.Close()

		return &buf, nil
	default:
		return data, nil
	}
}
