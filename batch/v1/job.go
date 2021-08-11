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
	"kmodules.xyz/resource-metrics/api"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
)

func init() {
	api.Register(batch.SchemeGroupVersion.WithKind("Job"), Job{}.ResourceCalculator())
}

type Job struct{}

func (r Job) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.DefaultPodRole},
		RuntimeRoles:           []api.PodRole{api.DefaultPodRole},
		RoleReplicasFn:         r.roleReplicasFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

func (_ Job) roleReplicasFn(obj map[string]interface{}) (api.ReplicaList, error) {
	return nil, nil
}

func (r Job) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
	return func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
		containers, err := api.AggregateContainerResources(obj, fn, api.AddResourceList, "spec", "template", "spec", "containers")
		if err != nil {
			return nil, err
		}
		initContainers, err := api.AggregateContainerResources(obj, fn, api.MaxResourceList, "spec", "template", "spec", "initContainers")
		if err != nil {
			return nil, err
		}
		return map[api.PodRole]core.ResourceList{
			api.DefaultPodRole: containers,
			api.InitPodRole:    initContainers,
		}, nil
	}
}
