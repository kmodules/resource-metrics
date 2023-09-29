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

package resourcemetrics

import (
	"fmt"
	"strings"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/management/v1alpha1"
	rsapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metrics/api"
	_ "kmodules.xyz/resource-metrics/apps/v1"
	_ "kmodules.xyz/resource-metrics/batch/v1"
	_ "kmodules.xyz/resource-metrics/batch/v1beta1"
	_ "kmodules.xyz/resource-metrics/core/v1"
	_ "kmodules.xyz/resource-metrics/kubedb.com/v1alpha2"
	_ "kmodules.xyz/resource-metrics/kubevault.com/v1alpha2"

	core "k8s.io/api/core/v1"
)

func Replicas(obj map[string]interface{}) (int64, error) {
	c, err := api.Load(obj)
	if err != nil {
		return 0, err
	}
	return c.Replicas(obj)
}

func RoleReplicas(obj map[string]interface{}) (api.ReplicaList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.RoleReplicas(obj)
}

func Mode(obj map[string]interface{}) (string, error) {
	c, err := api.Load(obj)
	if err != nil {
		return "", err
	}
	return c.Mode(obj)
}

func UsesTLS(obj map[string]interface{}) (bool, error) {
	c, err := api.Load(obj)
	if err != nil {
		return false, err
	}
	return c.UsesTLS(obj)
}

func TotalResourceLimits(obj map[string]interface{}) (core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.TotalResourceLimits(obj)
}

func TotalResourceRequests(obj map[string]interface{}) (core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.TotalResourceRequests(obj)
}

func AppResourceLimits(obj map[string]interface{}) (core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.AppResourceLimits(obj)
}

func AppResourceRequests(obj map[string]interface{}) (core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.AppResourceRequests(obj)
}

func RoleResourceLimits(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.RoleResourceLimits(obj)
}

func RoleResourceRequests(obj map[string]interface{}) (map[api.PodRole]core.ResourceList, error) {
	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}
	return c.RoleResourceRequests(obj)
}

func Quota(obj map[string]interface{}, pq *v1alpha1.ProjectQuota, apiType *kmapi.ResourceID) (*rsapi.QuotaDecision, error) {
	qd := &rsapi.QuotaDecision{
		Decision:   rsapi.DecisionAllow,
		Violations: make([]string, 0),
	}

	if pq == nil {
		qd.Decision = rsapi.DecisionNoOpinion
		return qd, nil
	}

	c, err := api.Load(obj)
	if err != nil {
		return nil, err
	}

	requests, err := c.AppResourceRequests(obj)
	if err != nil {
		return nil, err
	}
	limits, err := c.AppResourceLimits(obj)
	if err != nil {
		return nil, err
	}

	for _, quota := range pq.Status.Quotas {
		if quota.Group == apiType.Group {
			if quota.Kind != "" && quota.Kind != apiType.Kind {
				continue
			}
			hardRequests, hardLimits := extractRequestsLimits(quota.Hard)
			usedRequests, usedLimits := extractRequestsLimits(quota.Used)

			totRequestsUsage := api.AddResourceList(requests, usedRequests)
			for rn, usageQuan := range totRequestsUsage {
				hr, found := hardRequests[rn]
				if !found {
					continue
				}
				if usageQuan.Cmp(hr) > 0 {
					r := requests[rn]
					u := usedRequests[rn]
					l := hardRequests[rn]

					qd.Decision = rsapi.DecisionDeny
					qd.Violations = append(qd.Violations,
						fmt.Sprintf("Project quota exceeded. Requested: requests.%s=%s, Used: requests.%s=%s, Limited: requests.%s=%s", rn, r.String(), rn, u.String(), rn, l.String()))
				}
			}

			totLimitsUsage := api.AddResourceList(limits, usedLimits)
			for rn, usageQuan := range totLimitsUsage {
				hl, found := hardLimits[rn]
				if !found {
					continue
				}
				if usageQuan.Cmp(hl) > 0 {
					r := limits[rn]
					u := usedLimits[rn]
					l := hardLimits[rn]

					qd.Decision = rsapi.DecisionDeny
					qd.Violations = append(qd.Violations,
						fmt.Sprintf("Project quota exceeded. Requested: limits.%s=%s, Used: limits.%s=%s, Limited: limits.%s=%s", rn, r.String(), rn, u.String(), rn, l.String()))
				}
			}
		}
	}

	return qd, nil
}

func extractRequestsLimits(res core.ResourceList) (core.ResourceList, core.ResourceList) {
	requests := core.ResourceList{}
	limits := core.ResourceList{}

	for fullName, quan := range res {
		identifier, name, found := strings.Cut(fullName.String(), ".")
		if !found {
			continue
		}

		if identifier == "requests" {
			requests[core.ResourceName(name)] = quan
		} else {
			limits[core.ResourceName(name)] = quan
		}
	}

	return requests, limits
}
