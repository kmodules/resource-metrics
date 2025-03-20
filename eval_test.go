/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resourcemetrics

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tl "gomodules.xyz/testing"
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
			obj, err := tl.LoadFile(tt.name)
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
