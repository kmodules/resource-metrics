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
	"fmt"

	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	api.Register(schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha2",
		Kind:    "DocumentDB",
	}, DocumentDB{}.ResourceCalculator())
}

type DocumentDB struct{}

func (d DocumentDB) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRoleDefault},
		RuntimeRoles:           []api.PodRole{api.PodRoleDefault},
		RoleReplicasFn:         d.roleReplicasFn,
		ModeFn:                 d.modeFn,
		RoleResourceLimitsFn:   d.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: d.roleResourceFn(api.ResourceRequests),
	}
}

func (d DocumentDB) roleReplicasFn(obj map[string]any) (api.ReplicaList, error) {
	replicas, found, err := unstructured.NestedInt64(obj, "spec", "replicas")
	if err != nil {
		return nil, fmt.Errorf("failed to read spec.replicas %v: %w", obj, err)
	}
	if !found {
		replicas = 1
	}

	result := api.ReplicaList{
		api.PodRoleDefault: replicas,
	}
	return result, nil
}

func (d DocumentDB) modeFn(obj map[string]any) (string, error) {
	return DBModeStandalone, nil
}

func (d DocumentDB) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
	return func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
		container, replicas, err := api.AppNodeResourcesV2(obj, fn, DocumentDBContainerName, "spec")
		if err != nil {
			return nil, err
		}

		result := map[api.PodRole]api.PodInfo{
			api.PodRoleDefault: {Resource: container, Replicas: replicas},
		}
		return result, nil
	}
}
