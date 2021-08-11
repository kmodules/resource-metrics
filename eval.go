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
	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
)

// ExpressionFunction matches signature for https://github.com/gomodules/eval
type ExpressionFunction func(arguments ...interface{}) (interface{}, error)

func EvalFuncs() map[string]ExpressionFunction {
	return map[string]ExpressionFunction{
		"resource_replicas":       resource_replicas,
		"resource_mode":           resource_mode,
		"total_resource_limits":   total_resource_limits,
		"total_resource_requests": total_resource_requests,
		"app_resource_limits":     app_resource_limits,
		"app_resource_requests":   app_resource_requests,
	}
}

// resource_replicas(resource_obj)
func resource_replicas(args ...interface{}) (interface{}, error) {
	return api.Replicas(args[0].(map[string]interface{}))
}

// resource_mode(resource_obj)
func resource_mode(args ...interface{}) (interface{}, error) {
	return api.Mode(args[0].(map[string]interface{}))
}

// total_resource_limits(resource_obj, resource_type) => cpu cores (float64)
func total_resource_limits(args ...interface{}) (interface{}, error) {
	rr, err := api.TotalResourceLimits(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// total_resource_requests(resource_obj, resource_type)
func total_resource_requests(args ...interface{}) (interface{}, error) {
	rr, err := api.TotalResourceRequests(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// app_resource_limits(resource_obj, resource_type)
func app_resource_limits(args ...interface{}) (interface{}, error) {
	rr, err := api.AppResourceLimits(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// app_resource_requests(resource_obj, resource_type)
func app_resource_requests(args ...interface{}) (interface{}, error) {
	rr, err := api.AppResourceRequests(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}
