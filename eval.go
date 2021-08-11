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
	core "k8s.io/api/core/v1"
)

// EvalFuncs for https://github.com/gomodules/eval
func EvalFuncs() map[string]func(arguments ...interface{}) (interface{}, error) {
	return map[string]func(arguments ...interface{}) (interface{}, error){
		"resourceReplicas":      resourceReplicas,
		"resourceMode":          resourceMode,
		"totalResourceLimits":   totalResourceLimits,
		"totalResourceRequests": totalResourceRequests,
		"appResourceLimits":     appResourceLimits,
		"appResourceRequests":   appResourceRequests,
	}
}

// resourceReplicas(resource_obj)
func resourceReplicas(args ...interface{}) (interface{}, error) {
	return Replicas(args[0].(map[string]interface{}))
}

// resourceMode(resource_obj)
func resourceMode(args ...interface{}) (interface{}, error) {
	return Mode(args[0].(map[string]interface{}))
}

// totalResourceLimits(resource_obj, resource_type) => cpu cores (float64)
func totalResourceLimits(args ...interface{}) (interface{}, error) {
	rr, err := TotalResourceLimits(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// totalResourceRequests(resource_obj, resource_type)
func totalResourceRequests(args ...interface{}) (interface{}, error) {
	rr, err := TotalResourceRequests(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// appResourceLimits(resource_obj, resource_type)
func appResourceLimits(args ...interface{}) (interface{}, error) {
	rr, err := AppResourceLimits(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}

// appResourceRequests(resource_obj, resource_type)
func appResourceRequests(args ...interface{}) (interface{}, error) {
	rr, err := AppResourceRequests(args[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	q := rr[core.ResourceName(args[1].(string))]
	return q.AsApproximateFloat64(), nil
}
