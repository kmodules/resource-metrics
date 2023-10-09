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
	RegisterToPathMapperPlugin(&PostgresOpsRequest{})
}

type PostgresOpsRequest struct{}

var _ OpsPathMapper = (*PostgresOpsRequest)(nil)

func (m *PostgresOpsRequest) HorizontalPathMapping() map[OpsReqPath]ReferencedObjPath {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.horizontalScaling.replicas":      "",
		"spec.horizontalScaling.standbyMode":   "",
		"spec.horizontalScaling.streamingMode": "",
	}
}

func (m *PostgresOpsRequest) VerticalPathMapping() map[OpsReqPath]ReferencedObjPath {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.verticalScaling.postgres":    "",
		"spec.verticalScaling.exporter":    "",
		"spec.verticalScaling.coordinator": "",
	}
}

func (m *PostgresOpsRequest) VolumeExpansionPathMapping() map[OpsReqPath]ReferencedObjPath {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.volumeExpansion.postgres": "",
		"spec.volumeExpansion.mode":     "",
	}
}

func (m *PostgresOpsRequest) GetGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "PostgresOpsRequest",
	}
}
