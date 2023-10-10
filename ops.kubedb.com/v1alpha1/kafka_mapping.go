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
	RegisterToPathMapperPlugin(&kafkaOpsRequest{})
}

type kafkaOpsRequest struct{}

var _ OpsPathMapper = (*kafkaOpsRequest)(nil)

func (m *kafkaOpsRequest) HorizontalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.horizontalScaling.node":                "spec.replicas",
		"spec.horizontalScaling.topology.broker":     "spec.topology.broker.replicas",
		"spec.horizontalScaling.topology.controller": "spec.topology.controller.replicas",
	}, nil
}

func (m *kafkaOpsRequest) VerticalPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.verticalScaling.node":                "spec.podTemplate.spec.resources",
		"spec.verticalScaling.topology.broker":     "spec.topology.broker.resources",
		"spec.verticalScaling.topology.controller": "spec.topology.controller.resources",
	}, nil
}

func (m *kafkaOpsRequest) VolumeExpansionPathMapping(_ DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.volumeExpansion.node":                "spec.storage.resources.requests",
		"spec.volumeExpansion.topology.broker":     "spec.topology.broker.storage.resources.requests",
		"spec.volumeExpansion.topology.controller": "spec.topology.controller.storage.resources.requests",
	}, nil
}

func (m *kafkaOpsRequest) GetGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "KafkaOpsRequest",
	}
}
