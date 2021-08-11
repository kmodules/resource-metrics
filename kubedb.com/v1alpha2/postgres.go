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
		Kind:    "Postgres",
	}, Postgres{}.ResourceCalculator())
}

type Postgres struct{}

func (r Postgres) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.DefaultPodRole},
		RuntimeRoles:           []api.PodRole{api.DefaultPodRole, api.ExporterPodRole},
		RoleReplicasFn:         r.roleReplicasFn,
		ModeFn:                 r.modeFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

func (r Postgres) roleReplicasFn(obj map[string]interface{}) (api.ReplicaList, error) {
	v, found, err := unstructured.NestedInt64(obj, "spec", "replicas")
	if err != nil {
		return nil, fmt.Errorf("failed to read spec.replicas %v: %w", obj, err)
	}
	if !found {
		return api.ReplicaList{api.DefaultPodRole: 1}, nil
	}
	return api.ReplicaList{api.DefaultPodRole: v}, nil
}

func (r Postgres) modeFn(obj map[string]interface{}) (string, error) {
	mode, found, err := unstructured.NestedString(obj, "spec", "standbyMode")
	if err != nil {
		return "", err
	}
	if found && mode != "" {
		return mode, nil
	}
	return DBStandalone, nil
}

func (r Postgres) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
	return func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
		container, replicas, err := api.AppNodeResources(obj, fn, "spec")
		if err != nil {
			return nil, err
		}

		exporter, err := api.ContainerResources(obj, fn, "spec", "monitor", "prometheus", "exporter")
		if err != nil {
			return nil, err
		}
		return map[api.PodRole]core.ResourceList{
			api.DefaultPodRole:  api.MulResourceList(container, replicas),
			api.ExporterPodRole: api.MulResourceList(exporter, replicas),
		}, nil
	}
}
