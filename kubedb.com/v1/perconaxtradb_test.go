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

package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tl "gomodules.xyz/testing"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestPerconaXtraDB(t *testing.T) {
	type want struct {
		replicas       int64
		mode           string
		totalResources core.ResourceRequirements
		appResources   core.ResourceRequirements
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "testdata/kubedb.com/v1/percona-xtradb/standalone.yaml",
			want: want{
				replicas: 1,
				mode:     DBModeStandalone,
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("1000m"),
						core.ResourceMemory:  resource.MustParse("1128Mi"),
						core.ResourceStorage: resource.MustParse("1Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("500m"),
						core.ResourceMemory:  resource.MustParse("564Mi"),
						core.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
				appResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("500m"),
						core.ResourceMemory:  resource.MustParse("1000Mi"),
						core.ResourceStorage: resource.MustParse("1Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("250m"),
						core.ResourceMemory:  resource.MustParse("500Mi"),
						core.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		},
		{
			name: "testdata/kubedb.com/v1/percona-xtradb/cluster.yaml",
			want: want{
				replicas: 3,
				mode:     "Cluster",
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("3000m"),
						core.ResourceMemory:  resource.MustParse("4500Mi"),
						core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("1950m"),
						core.ResourceMemory:  resource.MustParse("2700Mi"),
						core.ResourceStorage: resource.MustParse("3Gi"),
					},
				},
				appResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("1500m"),
						core.ResourceMemory:  resource.MustParse("3000Mi"),
						core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:     resource.MustParse("750m"),
						core.ResourceMemory:  resource.MustParse("1500Mi"),
						core.ResourceStorage: resource.MustParse("3Gi"),
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
			c := PerconaXtraDB{}.ResourceCalculator()

			if got, err := c.Replicas(obj); err != nil {
				t.Errorf("Replicas() error = %v", err)
			} else if got != tt.want.replicas {
				t.Errorf("Replicas found = %v, expected = %v", got, tt.want.replicas)
			}

			if got, err := c.Mode(obj); err != nil {
				t.Errorf("Mode() error = %v", err)
			} else if got != tt.want.mode {
				t.Errorf("Mode found = %v, expected = %v", got, tt.want.mode)
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

			if got, err := c.AppResourceLimits(obj); err != nil {
				t.Errorf("AppResourceLimits() error = %v", err)
			} else if !cmp.Equal(tt.want.appResources.Limits, got) {
				t.Errorf("AppResourceLimits() difference = %v", cmp.Diff(tt.want.appResources.Limits, got))
			}
			if got, err := c.AppResourceRequests(obj); err != nil {
				t.Errorf("AppResourceRequests() error = %v", err)
			} else if !cmp.Equal(tt.want.appResources.Requests, got) {
				t.Errorf("AppResourceRequests() difference = %v", cmp.Diff(tt.want.appResources.Requests, got))
			}
		})
	}
}
