package gateway

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReadAndRestoreBody(t *testing.T) {
	tests := []struct {
		name        string
		input       io.ReadCloser
		expectedErr bool
		wantData    string
	}{
		{
			name:        "nil reader",
			input:       nil,
			expectedErr: false,
			wantData:    "",
		},
		{
			name:        "empty reader",
			input:       io.NopCloser(strings.NewReader("")),
			expectedErr: false,
			wantData:    "",
		},
		{
			name:        "successful read",
			input:       io.NopCloser(strings.NewReader("test data")),
			expectedErr: false,
			wantData:    "test data",
		},
		{
			name:        "error reading",
			input:       io.NopCloser(&errReader{}),
			expectedErr: true,
			wantData:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			restored, err := readAndRestoreBody(tt.input, buf)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if restored != nil {
					t.Error("expected nil reader on error but got non-nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.input == nil {
				if restored != nil {
					t.Error("expected nil reader but got non-nil")
				}
				return
			}

			data, err := io.ReadAll(restored)
			if err != nil {
				t.Errorf("error reading restored body: %v", err)
			}

			if got := string(data); got != tt.wantData {
				t.Errorf("restored data = %q, want %q", got, tt.wantData)
			}

			if got := buf.String(); got != tt.wantData {
				t.Errorf("buffer data = %q, want %q", got, tt.wantData)
			}
		})
	}
}

type errReader struct{}

func (r *errReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
