package util

const (
	URL_PREFIX                          = "http://"
	URL_HYPERCLUSTERRESOURCE_PATH       = "/hyperclusterresource"
	CONFIGMAP_NAME                      = "hypercloud-multi-agent-agentconfig"
	CONFIGMAP_NAMESPACE                 = "hypercloud-multi-agent-system"
	CONFIGMAP_MGNT_IP                   = "mgnt-ip"
	CONFIGMAP_MGNT_PORT                 = "mgnt-port"
	CONFIGMAP_MY_CLUSTERNAME            = "cluster-name"
	CONFIGMAP_REQUESTPERIOD             = "request-period"
	CONFIGMAP_RESOURCELIST              = "resourcelist"
	PROMETHEUS_SERVICE_NAME             = "prometheus-k8s"
	PROMETHEUS_SERVICE_NAMESPACE        = "monitoring"
	PROMETHEUS_NODE_EXPORT_SERVICE_NAME = "node-exporter"
	PROMETHEUS_QUERY_PATH               = "/api/v1/query"
	PROMETHEUS_QUERY_KEY_QUERY          = "query"
	PROMETHEUS_QUERY_KEY_TIME           = "time"
	PROMETHEUS_QUERY_POD_USAGE          = "count(kube_pod_info{host_ip=\"y.y.y.y\"})"
	PROMETHEUS_QUERY_CPU_USAGE          = "(   (1 - rate(node_cpu_seconds_total{job=\"node-exporter\", mode=\"idle\", instance=\"x.x.x.x:xxxx\"}[75s])) / ignoring(cpu) group_left   count without (cpu)( node_cpu_seconds_total{job=\"node-exporter\", mode=\"idle\", instance=\"x.x.x.x:xxxx\"}) )"
	PROMETHEUS_QUERY_STORAGE_USAGE      = "sum(   max by (device) (     node_filesystem_size_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\", fstype!=\"\"}   -     node_filesystem_avail_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\", fstype!=\"\"}   ) )"
	PROMETHEUS_QUERY_MEMORY_USAGE       = "((   node_memory_MemTotal_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_MemFree_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_Buffers_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_Cached_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} )) "
	LABEL_MASTER_ROLE                   = "node-role.kubernetes.io/master"
)
