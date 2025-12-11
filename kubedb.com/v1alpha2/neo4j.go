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
		Kind:    "Neo4j",
	}, Neo4j{}.ResourceCalculator())
}

type Neo4j struct{}

func (n Neo4j) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRoleDefault},
		RuntimeRoles:           []api.PodRole{api.PodRoleDefault, api.PodRoleExporter},
		RoleReplicasFn:         n.roleReplicasFn,
		ModeFn:                 n.modefn,
		UsesTLSFn:              n.usesTLSFn,
		RoleResourceLimitsFn:   n.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: n.roleResourceFn(api.ResourceRequests),
	}
}

func (n Neo4j) usesTLSFn(obj map[string]any) (bool, error) {
	_, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "tls")
	return found, err
}

func (n Neo4j) modefn(obj map[string]any) (string, error) {
	replicas, _, err := unstructured.NestedInt64(obj, "spec", "replicas")
	if err != nil {
		return "", err
	}
	if replicas > 1 {
		return DBModeCluster, nil
	}
	return DBModeStandalone, nil
}

func (n Neo4j) roleReplicasFn(obj map[string]any) (api.ReplicaList, error) {
	replicas, found, err := unstructured.NestedInt64(obj, "spec", "replicas")
	if err != nil {
		return nil, fmt.Errorf("failed to read spec.replicas %v: %w", obj, err)
	}
	if !found {
		return api.ReplicaList{api.PodRoleDefault: 1}, nil
	}
	return api.ReplicaList{api.PodRoleDefault: replicas}, nil
}

func (n Neo4j) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
	return func(obj map[string]any) (map[api.PodRole]api.PodInfo, error) {
		container, replicas, err := api.AppNodeResourcesV2(obj, fn, Neo4jContainerName, "spec")
		if err != nil {
			return nil, err
		}

		exporter, err := api.ContainerResources(obj, fn, "spec", "monitor", "prometheus", "exporter")
		if err != nil {
			return nil, err
		}
		return map[api.PodRole]api.PodInfo{
			api.PodRoleDefault:  {Resource: container, Replicas: replicas},
			api.PodRoleExporter: {Resource: exporter, Replicas: replicas},
		}, nil
	}
}
