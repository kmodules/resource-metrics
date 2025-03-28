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

package v1

const (
	DBModeEnsemble   = "Ensemble"
	DBModeCluster    = "Cluster"
	DBModeSentinel   = "Sentinel"
	DBModeSharded    = "Sharded"
	DBModeStandalone = "Standalone"
	DBModeReplicaSet = "ReplicaSet"
	DBModeDedicated  = "Dedicated"
	DBModeCombined   = "Combined"
)

const (
	ElasticsearchContainerName = "elasticsearch"
	MongoDBContainerName       = "mongodb"
	MySQLContainerName         = "mysql"
	PerconaXtraDBContainerName = "perconaxtradb"
	MariaDBContainerName       = "mariadb"
	PostgresContainerName      = "postgres"
	ProxySQLContainerName      = "proxysql"
	RedisContainerName         = "redis"
	PgBouncerContainerName     = "pgbouncer"
	KafkaContainerName         = "kafka"
	MemcachedContainerName     = "memcached"
	RedisSentinelContainerName = "redissentinel"
	MySQLRouterContainerName   = "mysql-router"

	// TODO: Update values
	MongoDBSidecarContainerName       = "replication-mode-detector"
	MySQLSidecarContainerName         = "mysql-coordinator"
	PerconaXtraDBSidecarContainerName = "px-coordinator"
	MariaDBSidecarContainerName       = "md-coordinator"
	PostgresSidecarContainerName      = "pg-coordinator"
	RedisSidecarContainerName         = "rd-coordinator"
	PgBouncerSidecarContainerName     = "pgbouncer"
	MySQLRouterSidecarContainerName   = "mysql-router"
)
