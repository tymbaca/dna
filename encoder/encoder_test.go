package encoder

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewDNAEncoder(t *testing.T) {
	tests := []struct {
		input   string
		want    []byte
		wantErr bool
	}{
		{
			input:   "CGGCATTA",
			want:    []byte{0b00011011},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.input)
			e := NewDNAEncoder(buf)

			resBuf := bytes.NewBuffer(nil)

			_, err := e.WriteTo(resBuf)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if !reflect.DeepEqual(resBuf.Bytes(), tt.want) {
				t.Errorf("incorrect result, want: %#v; got: %#v", tt.want, resBuf)
			}
		})
	}
}
