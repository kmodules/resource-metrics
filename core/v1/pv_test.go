package v1

import (
	"testing"
)

func TestPersistentVolume(t *testing.T) {
	type args struct {
		in0 map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := PersistentVolume{}
			if got := pe.Replicas(tt.args.in0); got != tt.want {
				t.Errorf("Replicas() = %v, want %v", got, tt.want)
			}
		})
	}
}
