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
	"reflect"

	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	api.Register(schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha2",
		Kind:    "ClickHouse",
	}, Pgpool{}.ResourceCalculator())
}

type ClickHouse struct{}

func (r ClickHouse) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRoleDefault},
		RuntimeRoles:           []api.PodRole{api.PodRoleDefault},
		RoleReplicasFn:         r.roleReplicasFn,
		ModeFn:                 r.modeFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

func (r ClickHouse) roleReplicasFn(obj map[string]interface{}) (api.ReplicaList, error) {
	result := api.ReplicaList{}

	clusterTopology, found, err := unstructured.NestedMap(obj, "spec", "clusterTopology")
	if err != nil {
		return nil, err
	}
	if found && clusterTopology != nil {
		// dedicated topology mode
		var replicas int64 = 0
		var shards int64 = 0

		clusters, _, err := unstructured.NestedSlice(clusterTopology, "cluster")
		if err != nil {
			return nil, err
		}

		for _, cluster := range clusters {
			shardCount, _, err := unstructured.NestedInt64(cluster.(map[string]interface{}), "shards")
			if err != nil {
				return nil, err
			}
			shards += shardCount
			shardReplicas, _, err := unstructured.NestedInt64(cluster.(map[string]interface{}), "replicas")
			if err != nil {
				return nil, err
			}
			replicas += shardReplicas * shards
		}
		result[api.PodRoleTotalShard] = replicas
		result[api.PodRoleShard] = shards

	} else {
		// standalone
		replicas, found, err := unstructured.NestedInt64(obj, "spec", "replicas")
		if err != nil {
			return nil, fmt.Errorf("failed to read spec.replicas %v: %w", obj, err)
		}
		if !found {
			result[api.PodRoleDefault] = 1
		} else {
			result[api.PodRoleDefault] = replicas
		}
	}
	return result, nil
}

func (r ClickHouse) modeFn(obj map[string]interface{}) (string, error) {
	clusterTopology, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "clusterTopology")
	if err != nil {
		return "", err
	}
	if found && !reflect.ValueOf(clusterTopology).IsNil() {
		return DBModeCluster, nil
	}
	return DBModeStandalone, nil
}

func (r ClickHouse) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]interface{}) (map[api.PodRole]api.PodInfo, error) {
	return func(obj map[string]interface{}) (map[api.PodRole]api.PodInfo, error) {
		container, replicas, err := api.AppNodeResourcesV2(obj, fn, ClickHouseContainerName, "spec")
		if err != nil {
			return nil, err
		}

		return map[api.PodRole]api.PodInfo{
			api.PodRoleDefault: {Resource: container, Replicas: replicas},
		}, nil
	}
}
