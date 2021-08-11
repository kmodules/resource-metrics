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

package api

import (
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type NotRegistered struct {
	gvk schema.GroupVersionKind
}

var _ error = NotRegistered{}

func (e NotRegistered) Error() string {
	return fmt.Sprintf("no calculator registered for %v", e.gvk)
}

func Replicas(obj map[string]interface{}) (int64, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return 0, NotRegistered{gvk}
	}
	return c.Replicas(obj)
}

func RoleReplicas(obj map[string]interface{}) (ReplicaList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.RoleReplicas(obj)
}

func Mode(obj map[string]interface{}) (string, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return "", NotRegistered{gvk}
	}
	return c.Mode(obj)
}

func TotalResourceLimits(obj map[string]interface{}) (core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.TotalResourceLimits(obj)
}

func TotalResourceRequests(obj map[string]interface{}) (core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.TotalResourceRequests(obj)
}

func AppResourceLimits(obj map[string]interface{}) (core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.AppResourceLimits(obj)
}

func AppResourceRequests(obj map[string]interface{}) (core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.AppResourceRequests(obj)
}

func RoleResourceLimits(obj map[string]interface{}) (map[PodRole]core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.RoleResourceLimits(obj)
}

func RoleResourceRequests(obj map[string]interface{}) (map[PodRole]core.ResourceList, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	c, ok := plugins[gvk]
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c.RoleResourceRequests(obj)
}
