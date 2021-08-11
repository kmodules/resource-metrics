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

	mt "kmodules.xyz/resource-metrics/testing"

	"github.com/google/go-cmp/cmp"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestPostgres(t *testing.T) {
	type want struct {
		replicas       int64
		totalResources core.ResourceRequirements
		appResources   core.ResourceRequirements
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "testdata/kubedb.com/v1beta2/postgres.yaml",
			want: want{
				replicas: 3,
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("1500m"),
						core.ResourceMemory: resource.MustParse("384Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("750m"),
						core.ResourceMemory: resource.MustParse("192Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
				},
				appResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("1500m"),
						core.ResourceMemory: resource.MustParse("384Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("750m"),
						core.ResourceMemory: resource.MustParse("192Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := mt.Load(tt.name)
			if err != nil {
				t.Error(err)
				return
			}
			c := Postgres{}.ResourceCalculator()

			if got, err := c.Replicas(obj); err != nil {
				t.Errorf("Replicas() error = %v", err)
			} else if got != tt.want.replicas {
				t.Errorf("Replicas = %v, want %v", got, tt.want)
			}

			if got, err := c.TotalResourceLimits(obj); err != nil {
				t.Errorf("TotalResourceLimits() error = %v", err)
			} else if !cmp.Equal(got, tt.want.totalResources.Limits) {
				t.Errorf("TotalResourceLimits() difference = %v", cmp.Diff(got, tt.want.totalResources.Limits))
			}
			if got, err := c.TotalResourceRequests(obj); err != nil {
				t.Errorf("TotalResourceRequests() error = %v", err)
			} else if !cmp.Equal(got, tt.want.totalResources.Requests) {
				t.Errorf("TotalResourceRequests() difference = %v", cmp.Diff(got, tt.want.totalResources.Requests))
			}

			if got, err := c.AppResourceLimits(obj); err != nil {
				t.Errorf("AppResourceLimits() error = %v", err)
			} else if !cmp.Equal(got, tt.want.appResources.Limits) {
				t.Errorf("AppResourceLimits() difference = %v", cmp.Diff(got, tt.want.appResources.Limits))
			}
			if got, err := c.AppResourceRequests(obj); err != nil {
				t.Errorf("AppResourceRequests() error = %v", err)
			} else if !cmp.Equal(got, tt.want.appResources.Requests) {
				t.Errorf("AppResourceRequests() difference = %v", cmp.Diff(got, tt.want.appResources.Requests))
			}
		})
	}
}
