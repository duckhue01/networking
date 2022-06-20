package json

import (
	"reflect"
	"testing"

	seri "github.com/duckhue01/lang/go/net/arch/serialization"
)

func TestJSON_Encode(t *testing.T) {
	type args struct {
		in0 []*seri.Chore
	}
	tests := []struct {
		name string
		r    *JSON
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JSON{}
			r.Encode(tt.args.in0)
		})
	}
}

func TestJSON_Decode(t *testing.T) {
	tests := []struct {
		name       string
		r          *JSON
		wantChores []*seri.Chore
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JSON{}
			gotChores, err := r.Decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChores, tt.wantChores) {
				t.Errorf("JSON.Decode() = %v, want %v", gotChores, tt.wantChores)
			}
		})
	}
}
