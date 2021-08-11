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
