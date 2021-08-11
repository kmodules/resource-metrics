package v1

import (
	"testing"

	mt "kmodules.xyz/resource-metrics/testing"

	"github.com/google/go-cmp/cmp"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestJob(t *testing.T) {
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
			name: "testdata/batch/v1/job.yaml",
			want: want{
				replicas: 0,
				totalResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("500m"),
						core.ResourceMemory: resource.MustParse("128Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("250m"),
						core.ResourceMemory: resource.MustParse("64Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
				},
				appResources: core.ResourceRequirements{
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("500m"),
						core.ResourceMemory: resource.MustParse("128Mi"),
						// core.ResourceStorage: resource.MustParse("3Gi"),
					},
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("250m"),
						core.ResourceMemory: resource.MustParse("64Mi"),
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
			c := Job{}.ResourceCalculator()

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
