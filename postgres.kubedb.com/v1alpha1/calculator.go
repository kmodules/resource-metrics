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

package v1alpha1

import (
	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
)

// Publisher and Subscriber only configure logical replication on an existing Postgres
// instance; no pod is created for them, so every role reports zero.
func noResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRoleDefault},
		RuntimeRoles:           []api.PodRole{api.PodRoleDefault},
		RoleReplicasFn:         noRoleReplicas,
		RoleResourceLimitsFn:   noRoleResources,
		RoleResourceRequestsFn: noRoleResources,
	}
}

func noRoleReplicas(_ map[string]any) (api.ReplicaList, error) {
	return api.ReplicaList{api.PodRoleDefault: 0}, nil
}

func noRoleResources(_ map[string]any) (map[api.PodRole]api.PodInfo, error) {
	return map[api.PodRole]api.PodInfo{
		api.PodRoleDefault: {Resource: core.ResourceList{}, Replicas: 0},
	}, nil
}
