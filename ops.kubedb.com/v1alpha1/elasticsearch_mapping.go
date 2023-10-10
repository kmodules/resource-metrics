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
	RegisterToPathMapperPlugin(&ElasticsearchOpsRequest{})
}

type ElasticsearchOpsRequest struct{}

var _ OpsPathMapper = (*ElasticsearchOpsRequest)(nil)

func (m *ElasticsearchOpsRequest) HorizontalPathMapping(dbObj DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.horizontalScaling.node":                  "",
		"spec.horizontalScaling.topology.master":       "",
		"spec.horizontalScaling.topology.ingest":       "",
		"spec.horizontalScaling.topology.data":         "",
		"spec.horizontalScaling.topology.dataContent":  "",
		"spec.horizontalScaling.topology.dataHot":      "",
		"spec.horizontalScaling.topology.dataWarm":     "",
		"spec.horizontalScaling.topology.dataCold":     "",
		"spec.horizontalScaling.topology.dataFrozen":   "",
		"spec.horizontalScaling.topology.ml":           "",
		"spec.horizontalScaling.topology.transform":    "",
		"spec.horizontalScaling.topology.coordinating": "",
	}, nil
}

func (m *ElasticsearchOpsRequest) VerticalPathMapping(dbObj DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.verticalScaling.node":                  "",
		"spec.verticalScaling.exporter":              "",
		"spec.verticalScaling.topology.master":       "",
		"spec.verticalScaling.topology.ingest":       "",
		"spec.verticalScaling.topology.data":         "",
		"spec.verticalScaling.topology.dataContent":  "",
		"spec.verticalScaling.topology.dataHot":      "",
		"spec.verticalScaling.topology.dataWarm":     "",
		"spec.verticalScaling.topology.dataCold":     "",
		"spec.verticalScaling.topology.dataFrozen":   "",
		"spec.verticalScaling.topology.ml":           "",
		"spec.verticalScaling.topology.transform":    "",
		"spec.verticalScaling.topology.coordinating": "",
	}, nil
}

func (m *ElasticsearchOpsRequest) VolumeExpansionPathMapping(dbObj DbObject) (map[OpsReqPath]ReferencedObjPath, error) {
	return map[OpsReqPath]ReferencedObjPath{
		"spec.volumeExpansion.mode":                  "",
		"spec.volumeExpansion.node":                  "",
		"spec.volumeExpansion.topology.master":       "",
		"spec.volumeExpansion.topology.ingest":       "",
		"spec.volumeExpansion.topology.data":         "",
		"spec.volumeExpansion.topology.dataContent":  "",
		"spec.volumeExpansion.topology.dataHot":      "",
		"spec.volumeExpansion.topology.dataWarm":     "",
		"spec.volumeExpansion.topology.dataCold":     "",
		"spec.volumeExpansion.topology.dataFrozen":   "",
		"spec.volumeExpansion.topology.ml":           "",
		"spec.volumeExpansion.topology.transform":    "",
		"spec.volumeExpansion.topology.coordinating": "",
	}, nil
}

func (m *ElasticsearchOpsRequest) GetGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    "ElasticsearchOpsRequest",
	}
}
