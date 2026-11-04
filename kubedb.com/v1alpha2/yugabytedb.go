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
	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DefaultYugabyteDBReplicationFactor matches the operator's default (spec.replicationFactor): the
// number of masters, and the default number of tservers, of a universe that does not set one.
const DefaultYugabyteDBReplicationFactor int64 = 3

func init() {
	api.Register(schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha2",
		Kind:    "YugabyteDB",
	}, YugabyteDB{}.ResourceCalculator())
	api.Register(schema.GroupVersionKind{
		Group:   "gitops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "YugabyteDB",
	}, YugabyteDB{}.ResourceCalculator())
}

type YugabyteDB struct{}

func (r YugabyteDB) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRoleMaster, api.PodRoleTServer, api.PodRoleReadReplica},
		RuntimeRoles:           []api.PodRole{api.PodRoleMaster, api.PodRoleTServer, api.PodRoleReadReplica},
		RoleReplicasFn:         r.roleReplicasFn,
		ModeFn:                 r.modeFn,
		UsesTLSFn:              r.usesTLSFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

// yugabyteDBReplicationFactor reads spec.replicationFactor, defaulting as the operator does.
func (r YugabyteDB) yugabyteDBReplicationFactor(obj map[string]any) (int64, error) {
	rf, found, err := unstructured.NestedInt64(obj, "spec", "replicationFactor")
	if err != nil {
		return 0, err
	}
	if !found {
		return DefaultYugabyteDBReplicationFactor, nil
	}
	return rf, nil
}

// roleReplicasFn reads the master and tserver pod counts, and the read replica count when a read
// replica cluster is configured. The master count always equals the replication factor; the
// operator rejects any other value.
func (r YugabyteDB) roleReplicasFn(obj map[string]any) (api.ReplicaList, error) {
	rf, err := r.yugabyteDBReplicationFactor(obj)
	if err != nil {
		return nil, err
	}
	result := api.ReplicaList{api.PodRoleMaster: rf}

	tserverReplicas, found, err := unstructured.NestedInt64(obj, "spec", "topology", "tserver", "replicas")
	if err != nil {
		return nil, err
	}
	if !found {
		tserverReplicas = rf
	}
	result[api.PodRoleTServer] = tserverReplicas

	readReplica, found, err := unstructured.NestedMap(obj, "spec", "topology", "readReplica")
	if err != nil {
		return nil, err
	}
	if found && readReplica != nil {
		rrReplicas, found, err := unstructured.NestedInt64(readReplica, "replicas")
		if err != nil {
			return nil, err
		}
		if !found {
			rrReplicas = 1
		}
		result[api.PodRoleReadReplica] = rrReplicas
	}

	return result, nil
}

// modeFn is always Distributed: a universe always runs the master and tserver roles as separate
// pods, unlike the single process/combined mode some sibling databases also support.
func (r YugabyteDB) modeFn(_ map[string]any) (string, error) {
	return DBModeDistributed, nil
}

func (r YugabyteDB) usesTLSFn(obj map[string]any) (bool, error) {
	_, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "tls")
	return found, err
}

func (r YugabyteDB) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
	return func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
		master, masterReplicas, err := api.AppNodeResourcesV2(obj, fn, YugabyteDBMasterContainerName, "spec", "topology", "master")
		if err != nil {
			return nil, err
		}
		tserver, tserverReplicas, err := api.AppNodeResourcesV2(obj, fn, YugabyteDBTServerContainerName, "spec", "topology", "tserver")
		if err != nil {
			return nil, err
		}

		result := map[api.PodRole]api.PodInfo{
			api.PodRoleMaster:  {Resource: master, Replicas: masterReplicas},
			api.PodRoleTServer: {Resource: tserver, Replicas: tserverReplicas},
		}

		readReplica, found, err := unstructured.NestedMap(obj, "spec", "topology", "readReplica")
		if err != nil {
			return nil, err
		}
		if found && readReplica != nil {
			// The read replica tservers run the same container, inlined under readReplica.
			rr, rrReplicas, err := api.AppNodeResourcesV2(obj, fn, YugabyteDBTServerContainerName, "spec", "topology", "readReplica")
			if err != nil {
				return nil, err
			}
			result[api.PodRoleReadReplica] = api.PodInfo{Resource: rr, Replicas: rrReplicas}
		}

		return result, nil
	}
}
