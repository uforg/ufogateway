package gateway

import (
	"bytes"
	"io"
)

// readAndRestoreBody reads from the provided io.ReadCloser into the provided bytes.Buffer,
// then restores the io.ReadCloser so it can be read again.
// It returns the restored io.ReadCloser and any error encountered.
func readAndRestoreBody(rc io.ReadCloser, buf *bytes.Buffer) (io.ReadCloser, error) {
	if rc != nil {
		_, err := buf.ReadFrom(rc)
		if err != nil {
			return nil, err
		}
		// Restore the io.ReadCloser with the data read
		rc = io.NopCloser(bytes.NewReader(buf.Bytes()))
	}
	return rc, nil
}
