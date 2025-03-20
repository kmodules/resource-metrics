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
	"testing"

	"gomodules.xyz/encoding/json"
	tl "gomodules.xyz/testing"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type testSet struct {
	name string
	want want
}
type want map[ReferencedObjPath]interface{}

func TestMongoDBScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/horizontal/cluster.yaml",
			want: want{
				"spec.replicas": 4,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/horizontal/sharding.yaml",
			want: want{
				"spec.shardTopology.shard.shards":          3,
				"spec.shardTopology.shard.replicas":        4,
				"spec.shardTopology.configServer.replicas": 4,
				"spec.shardTopology.mongos.replicas":       3,
				"spec.hidden.replicas":                     2,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/vertical/cluster.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.6"),
						core.ResourceMemory: resource.MustParse("1.2Gi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.6"),
						core.ResourceMemory: resource.MustParse("1.2Gi"),
					},
				},
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/vertical/sharding.yaml",
			want: want{
				"spec.shardTopology.shard.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceMemory: resource.MustParse("1100Mi"),
						core.ResourceCPU:    resource.MustParse("0.55"),
					},
					Limits: core.ResourceList{
						core.ResourceMemory: resource.MustParse("1100Mi"),
						core.ResourceCPU:    resource.MustParse("0.55"),
					},
				},
				"spec.shardTopology.configServer.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.56"),
						core.ResourceMemory: resource.MustParse("1110Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.56"),
						core.ResourceMemory: resource.MustParse("1110Mi"),
					},
				},
				"spec.shardTopology.mongos.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.57"),
						core.ResourceMemory: resource.MustParse("1120Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.57"),
						core.ResourceMemory: resource.MustParse("1120Mi"),
					},
				},
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/vertical/standalone.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("1"),
						core.ResourceMemory: resource.MustParse("2Gi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("1"),
						core.ResourceMemory: resource.MustParse("2Gi"),
					},
				},
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/volumeexpansion/standalone.yaml",
			want: want{
				"spec.storage.resources.requests.storage": resource.MustParse("2Gi"),
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/volumeexpansion/sharding.yaml",
			want: want{
				"spec.shardTopology.shard.storage.resources.requests.storage":        resource.MustParse("3Gi"),
				"spec.shardTopology.configServer.storage.resources.requests.storage": resource.MustParse("4Gi"),
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mongodbops/volumeexpansion/cluster.yaml",
			want: want{
				"spec.storage.resources.requests.storage": resource.MustParse("5Gi"),
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestMySQLScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mysqlops/horizontal/cluster.yaml",
			want: want{
				"spec.replicas": 5,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mysqlops/vertical/standalone.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("700m"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.7"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
				},
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/mysqlops/vertical/cluster.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("700m"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.7"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
				},
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestProxySQLScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/proxysqlops/horizontal.yaml",
			want: want{
				"spec.replicas": 5,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/proxysqlops/vertical.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.6"),
						core.ResourceMemory: resource.MustParse("1.2Gi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.7"),
						core.ResourceMemory: resource.MustParse("1.3Gi"),
					},
				},
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestRedisScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/redisops/horizontal/cluster.yaml",
			want: want{
				"spec.cluster.master":   4,
				"spec.cluster.replicas": 1,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/redisops/vertical/standalone.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("200m"),
						core.ResourceMemory: resource.MustParse("300Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("500m"),
						core.ResourceMemory: resource.MustParse("800Mi"),
					},
				},
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/redisops/vertical/cluster.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("200m"),
						core.ResourceMemory: resource.MustParse("300Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("500m"),
						core.ResourceMemory: resource.MustParse("800Mi"),
					},
				},
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestRedisSentinelScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/redissentinel/horizontal/sentinel.yaml",
			want: want{
				"spec.replicas": 3,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/redissentinel/vertical/sentinel.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("200m"),
						core.ResourceMemory: resource.MustParse("300Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("500m"),
						core.ResourceMemory: resource.MustParse("800Mi"),
					},
				},
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestPostgresScaling(t *testing.T) {
	tests := []testSet{
		{
			name: "testdata/ops.kubedb.com/v1alpha1/postgresops/horizontal.yaml",
			want: want{
				"spec.replicas": 5,
			},
		},
		{
			name: "testdata/ops.kubedb.com/v1alpha1/postgresops/vertical.yaml",
			want: want{
				"spec.podTemplate.spec.resources": core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("700m"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
					Limits: core.ResourceList{
						core.ResourceCPU:    resource.MustParse("0.7"),
						core.ResourceMemory: resource.MustParse("1200Mi"),
					},
				},
			},
		},
	}

	evaluateTestSet(t, tests)
}

func TestKafkaScaling(t *testing.T) {
	tests := []testSet{}

	evaluateTestSet(t, tests)
}

func TestElasticSearchScaling(t *testing.T) {
	tests := []testSet{}

	evaluateTestSet(t, tests)
}

func TestMariaDBScaling(t *testing.T) {
	tests := []testSet{}

	evaluateTestSet(t, tests)
}

func TestMemcachedScaling(t *testing.T) {
	tests := []testSet{}

	evaluateTestSet(t, tests)
}

func TestPgBouncerScaling(t *testing.T) {
	tests := []testSet{}

	evaluateTestSet(t, tests)
}

func evaluateTestSet(t *testing.T, tests []testSet) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opsObj, err := tl.LoadFile(tt.name)
			if err != nil {
				t.Error(err)
				return
			}
			scaledObject, err := GetScaledObject(opsObj)
			if err != nil {
				t.Error(err)
				return
			}

			for k, ev := range tt.want {
				sop := splitPathToSlice(string(k))
				got, found, err := unstructured.NestedFieldNoCopy(scaledObject, sop...)
				if err != nil {
					t.Error(err)
					return
				}

				if !found || got == nil {
					t.Errorf("For path `%s`, found = `%v`, expected = `%v`", k, got, ev)
				} else {
					if !isEqual(ev, got) {
						t.Errorf("For path `%s`, found = `%v`, expected = `%v`", k, got, ev)
					}
				}
			}
		})
	}
}

func isEqual(x, y interface{}) bool {
	if checkIntType(x, y) {
		return compareIntType(x, y)
	} else if checkResourceRequirementsType(x, y) {
		return compareResourceRequirementsType(x, y)
	} else if checkResourceQuantityType(x, y) {
		return compareResourceQuantityType(x, y)
	}

	return false
}

func checkIntType(x, y interface{}) bool {
	_, okx := x.(int64)
	_, okx1 := x.(int)
	_, oky := y.(int64)
	_, oky1 := y.(int64)

	if okx || oky || okx1 || oky1 {
		return true
	}

	return false
}

func checkResourceRequirementsType(x, y interface{}) bool {
	_, okx := isResourceRequirementsType(x)
	_, oky := isResourceRequirementsType(y)
	return okx || oky
}

func checkResourceQuantityType(x, y interface{}) bool {
	_, okx := isResourceQuantityType(x)
	_, oky := isResourceQuantityType(y)
	return okx || oky
}

func isResourceQuantityType(x interface{}) (resource.Quantity, bool) {
	v, ok := x.(resource.Quantity)
	return v, ok
}

func isResourceRequirementsType(x interface{}) (core.ResourceRequirements, bool) {
	v, okx := x.(core.ResourceRequirements)
	return v, okx
}

func compareIntType(x, y interface{}) bool {
	xx, ok := x.(int64)
	if !ok {
		xx = int64(x.(int))
	}
	yy, ok := y.(int64)
	if !ok {
		yy = int64(y.(int))
	}

	return xx == yy
}

func compareResourceRequirementsType(x, y interface{}) bool {
	xx := core.ResourceRequirements{
		Requests: make(core.ResourceList),
		Limits:   make(core.ResourceList),
	}
	yy := core.ResourceRequirements{
		Requests: make(core.ResourceList),
		Limits:   make(core.ResourceList),
	}

	if v, ok := isMapStringInterfaceType(x); ok {
		xx = convertMapStringInterfaceToResourceRequirements(v)
	} else if v, ok = isMapStringInterfaceType(y); ok {
		yy = convertMapStringInterfaceToResourceRequirements(v)
	}

	if v, ok := isResourceRequirementsType(x); ok {
		xx = v
	} else if v, ok := isResourceRequirementsType(y); ok {
		yy = v
	}

	return compareResourceList(xx.Requests, yy.Requests) && compareResourceList(xx.Limits, yy.Limits)
}

func compareResourceList(x, y core.ResourceList) (ok bool) {
	ok = true
	for kx, vx := range x {
		vy, found := y[kx]
		if !found {
			ok = false
			return
		}
		if vx.Cmp(vy) != 0 {
			ok = false
			return
		}
	}

	return
}

func compareResourceQuantityType(x, y interface{}) bool {
	xx := resource.Quantity{}
	yy := resource.Quantity{}

	if v, ok := isResourceQuantityType(x); ok {
		xx = v
	} else if v, ok = isResourceQuantityType(y); ok {
		yy = v
	}

	if v, ok := x.(string); ok {
		xx = resource.MustParse(v)
	} else if v, ok := y.(string); ok {
		yy = resource.MustParse(v)
	}

	return xx.Equal(yy)
}

func isMapStringInterfaceType(x interface{}) (map[string]interface{}, bool) {
	v, ok := x.(map[string]interface{})
	return v, ok
}

func convertMapStringInterfaceToResourceRequirements(x map[string]interface{}) core.ResourceRequirements {
	k := core.ResourceRequirements{
		Requests: make(core.ResourceList),
		Limits:   make(core.ResourceList),
	}
	bytes, _ := json.Marshal(x)
	_ = json.Unmarshal(bytes, &k)

	return k
}
