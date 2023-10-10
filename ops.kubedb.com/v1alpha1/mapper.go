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
	"errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/resource-metrics/api"
)

type OpsPathMapper interface {
	HorizontalPathMapping(map[string]interface{}) map[OpsReqPath]ReferencedObjPath
	VerticalPathMapping(map[string]interface{}) map[OpsReqPath]ReferencedObjPath
	VolumeExpansionPathMapping(map[string]interface{}) map[OpsReqPath]ReferencedObjPath
	GetGroupVersionKind() schema.GroupVersionKind
}

type OpsReqPath string
type ReferencedObjPath string
type ScaledObject map[string]interface{}

var PathMapperPlugin map[schema.GroupVersionKind]OpsPathMapper

func RegisterToPathMapperPlugin(opsObj OpsPathMapper) {
	PathMapperPlugin[opsObj.GetGroupVersionKind()] = opsObj
}

func LoadOpsPathMapper(opsObj map[string]interface{}) (OpsPathMapper, error) {
	gvk := getGVK(opsObj)
	opsMapperObj, found := PathMapperPlugin[gvk]
	if !found {
		return nil, errors.New("gvk not registered")
	}

	return opsMapperObj, nil
}

func RegisterPathMapperPluginMembersWithApiPlugin(rc api.ResourceCalculator) {
	for _, pm := range PathMapperPlugin {
		api.Register(pm.GetGroupVersionKind(), rc)
	}
}
