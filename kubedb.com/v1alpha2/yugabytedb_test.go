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

package v1alpha2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tl "gomodules.xyz/testing"
	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestYugabyteDB(t *testing.T) {
	type want struct {
		replicas       int64
		mode           string
		usesTLS        bool
		totalResources core.ResourceRequirements
	}
	tests := []struct {
		name string
		want want
	}{
		{
			// master (RF3) + tserver (RF3), no read replica, TLS on.
			name: "testdata/kubedb.com/v1alpha2/yugabytedb/cluster.yaml",
			want: want{
				replicas: 6,
				mode:     DBModeDistributed,
				usesTLS:  true,
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("9"),
						core.ResourceMemory:  resource.MustParse("18Gi"),
						core.ResourceStorage: resource.MustParse("180Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("4500m"),
						core.ResourceMemory:  resource.MustParse("9Gi"),
						core.ResourceStorage: resource.MustParse("180Gi"),
					},
				},
			},
		},
		{
			// master (RF3) + tserver (RF3) + a 2 pod read replica, no TLS.
			name: "testdata/kubedb.com/v1alpha2/yugabytedb/with-read-replica.yaml",
			want: want{
				replicas: 8,
				mode:     DBModeDistributed,
				usesTLS:  false,
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("11"),
						core.ResourceMemory:  resource.MustParse("22Gi"),
						core.ResourceStorage: resource.MustParse("220Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("5500m"),
						core.ResourceMemory:  resource.MustParse("11Gi"),
						core.ResourceStorage: resource.MustParse("220Gi"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := tl.LoadFile(tt.name)
			if err != nil {
				t.Error(err)
				return
			}
			c := YugabyteDB{}.ResourceCalculator()

			if got, err := c.Replicas(obj); err != nil {
				t.Errorf("Replicas() error = %v", err)
			} else if got != tt.want.replicas {
				t.Errorf("Replicas() got = %v, want %v", got, tt.want.replicas)
			}

			if got, err := c.Mode(obj); err != nil {
				t.Errorf("Mode() error = %v", err)
			} else if got != tt.want.mode {
				t.Errorf("Mode() got = %v, want %v", got, tt.want.mode)
			}

			if got, err := c.UsesTLS(obj); err != nil {
				t.Errorf("UsesTLS() error = %v", err)
			} else if got != tt.want.usesTLS {
				t.Errorf("UsesTLS() got = %v, want %v", got, tt.want.usesTLS)
			}

			if got, err := c.TotalResourceLimits(obj); err != nil {
				t.Errorf("TotalResourceLimits() error = %v", err)
			} else if !cmp.Equal(tt.want.totalResources.Limits, got) {
				t.Errorf("TotalResourceLimits() difference = %v", cmp.Diff(tt.want.totalResources.Limits, got))
			}
			if got, err := c.TotalResourceRequests(obj); err != nil {
				t.Errorf("TotalResourceRequests() error = %v", err)
			} else if !cmp.Equal(tt.want.totalResources.Requests, got) {
				t.Errorf("TotalResourceRequests() difference = %v", cmp.Diff(tt.want.totalResources.Requests, got))
			}
		})
	}
}

// TestYugabyteDBRoleReplicas checks that each role's pod count is reported separately, not only
// the total: the master count always equals the replication factor.
func TestYugabyteDBRoleReplicas(t *testing.T) {
	obj, err := tl.LoadFile("testdata/kubedb.com/v1alpha2/yugabytedb/with-read-replica.yaml")
	if err != nil {
		t.Fatal(err)
	}
	c := YugabyteDB{}.ResourceCalculator()
	got, err := c.RoleReplicas(obj)
	if err != nil {
		t.Fatal(err)
	}
	want := api.ReplicaList{
		api.PodRoleMaster:      3,
		api.PodRoleTServer:     3,
		api.PodRoleReadReplica: 2,
	}
	if !cmp.Equal(want, got) {
		t.Errorf("RoleReplicas() difference = %v", cmp.Diff(want, got))
	}
}

// TestYugabyteDBDefaultReplicas checks the defaults a CR has before the operator's webhook runs:
// the master and tserver counts fall back to the replication factor.
func TestYugabyteDBDefaultReplicas(t *testing.T) {
	obj := map[string]any{
		"apiVersion": "kubedb.com/v1alpha2",
		"kind":       "YugabyteDB",
		"spec":       map[string]any{},
	}
	c := YugabyteDB{}.ResourceCalculator()
	got, err := c.RoleReplicas(obj)
	if err != nil {
		t.Fatal(err)
	}
	want := api.ReplicaList{
		api.PodRoleMaster:  DefaultYugabyteDBReplicationFactor,
		api.PodRoleTServer: DefaultYugabyteDBReplicationFactor,
	}
	if !cmp.Equal(want, got) {
		t.Errorf("RoleReplicas() difference = %v", cmp.Diff(want, got))
	}
}
