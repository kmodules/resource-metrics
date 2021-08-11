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

type PodRole string

const (
	DefaultPodRole      PodRole = ""
	InitPodRole         PodRole = "init"
	RouterPodRole       PodRole = "router"
	ExporterPodRole     PodRole = "exporter"
	TotalShardPodRole   PodRole = "total_shard"
	ShardPodRole        PodRole = "shard"
	PerShardPodRole     PodRole = "per_shard"
	ConfigServerPodRole PodRole = "config_server"
	MongosPodRole       PodRole = "mongos"
)

type ReplicaList map[PodRole]int64
