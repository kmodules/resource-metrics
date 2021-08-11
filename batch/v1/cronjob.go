package v1

import (
	"kmodules.xyz/resource-metrics/api"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
)

func init() {
	api.Register(batch.SchemeGroupVersion.WithKind("CronJob"), CronJob{}.ResourceCalculator())
}

type CronJob struct{}

func (r CronJob) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.DefaultPodRole},
		RuntimeRoles:           []api.PodRole{api.DefaultPodRole},
		RoleReplicasFn:         r.roleReplicasFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

func (_ CronJob) roleReplicasFn(obj map[string]interface{}) (api.ReplicaList, error) {
	return nil, nil
}

func (r CronJob) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
	return func(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
		containers, err := api.AggregateContainerResources(obj, fn, api.AddResourceList, "spec", "jobTemplate", "spec", "template", "spec", "containers")
		if err != nil {
			return nil, err
		}
		initContainers, err := api.AggregateContainerResources(obj, fn, api.MaxResourceList, "spec", "jobTemplate", "spec", "template", "spec", "initContainers")
		if err != nil {
			return nil, err
		}
		return map[api.PodRole]core.ResourceList{
			api.DefaultPodRole: containers,
			api.InitPodRole:    initContainers,
		}, nil
	}
}
