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

import "k8s.io/apimachinery/pkg/runtime/schema"

func init() {
	RegisterToPathMapperPlugin(&ProxySQLOpsRequest{})
}

type ProxySQLOpsRequest struct{}

var _ OpsPathMapper = (*ProxySQLOpsRequest)(nil)

func (m *ProxySQLOpsRequest) HorizontalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.horizontalScaling.member": "",
	}, nil
}

func (m *ProxySQLOpsRequest) VerticalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.verticalScaling.proxysql": "",
	}, nil
}

func (m *ProxySQLOpsRequest) VolumeExpansionPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{}, nil
}

func (m *ProxySQLOpsRequest) GetGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "ProxySQLOpsRequest",
	}
}
