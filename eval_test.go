package resourcemetrics

import (
	"testing"

	mt "kmodules.xyz/resource-metrics/testing"

	"github.com/google/go-cmp/cmp"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_totalResourceLimits(t *testing.T) {
	tests := []struct {
		name    string
		args    []interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "testdata/kubedb.com/v1alpha2/mongodb/replicaset.yaml",
			args: []interface{}{
				core.ResourceMemory,
			},
			want:    resourceQuantityAsFloat64(core.ResourceMemory, resource.MustParse("3Gi")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := mt.Load(tt.name)
			if err != nil {
				t.Error(err)
				return
			}
			got, err := totalResourceLimits(append([]interface{}{obj}, tt.args...)...)
			if (err != nil) != tt.wantErr {
				t.Errorf("totalResourceLimits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(tt.want, got) {
				t.Errorf("totalResourceLimits()difference = %v", cmp.Diff(tt.want, got))
			}
		})
	}
}
