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
	RegisterToPathMapperPlugin(&MySqlOpsRequest{})
}

type MySqlOpsRequest struct{}

var _ OpsPathMapper = (*MySqlOpsRequest)(nil)

func (m *MySqlOpsRequest) HorizontalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.horizontalScaling.member": "spec.replicas",
	}, nil
}

func (m *MySqlOpsRequest) VerticalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.verticalScaling.mysql":       "",
		"spec.verticalScaling.exporter":    "spec.monitor.prometheus.exporter.resources",
		"spec.verticalScaling.coordinator": "spec.coordinator.resources",
	}, nil
}

func (m *MySqlOpsRequest) VolumeExpansionPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.volumeExpansion.mysql": "",
		"spec.volumeExpansion.mode":  "",
	}, nil
}

func (m *MySqlOpsRequest) GetGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "MySQLOpsRequest",
	}
}