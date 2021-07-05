package helm

import (
	"crypto/tls"
	"net/http"
	"os/exec"
	"time"

	"k8s.io/klog"
	// "k8s.io/kubectl/pkg/cmd/annotate"
)

const (
	URL_PREFIX                              = "http://"
	URL_HYPERCLUSTERRESOURCE_PATH           = "/hyperclusterresource"
	URL_HTPERCLOUD_URL_PATH                 = "/hypercloudurl"
	URL_PARAM_URL                           = "url"
	CONFIGMAP_NAME                          = "hypercloud-multi-agent-agentconfig"
	CONFIGMAP_NAMESPACE                     = "hypercloud-multi-agent-system"
	CONFIGMAP_MGNT_IP                       = "mgnt-ip"
	CONFIGMAP_MGNT_PORT                     = "mgnt-port"
	CONFIGMAP_MY_CLUSTERNAME                = "cluster-name"
	CONFIGMAP_REQUESTPERIOD                 = "request-period"
	CONFIGMAP_RESOURCELIST                  = "resourcelist"
	PROMETHEUS_SERVICE_NAME                 = "prometheus-k8s"
	PROMETHEUS_SERVICE_NAMESPACE            = "monitoring"
	PROMETHEUS_NODE_EXPORT_SERVICE_NAME     = "node-exporter"
	PROMETHEUS_QUERY_PATH                   = "/api/v1/query"
	PROMETHEUS_QUERY_KEY_QUERY              = "query"
	PROMETHEUS_QUERY_KEY_TIME               = "time"
	PROMETHEUS_QUERY_POD_USAGE              = "count(kube_pod_info{host_ip=\"y.y.y.y\"})"
	PROMETHEUS_QUERY_CPU_USAGE              = "(   (1 - rate(node_cpu_seconds_total{job=\"node-exporter\", mode=\"idle\", instance=\"x.x.x.x:xxxx\"}[75s])) / ignoring(cpu) group_left   count without (cpu)( node_cpu_seconds_total{job=\"node-exporter\", mode=\"idle\", instance=\"x.x.x.x:xxxx\"}) )"
	PROMETHEUS_QUERY_STORAGE_USAGE          = "sum(   max by (device) (     node_filesystem_size_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\", fstype!=\"\"}   -     node_filesystem_avail_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\", fstype!=\"\"}   ) )"
	PROMETHEUS_QUERY_MEMORY_USAGE           = "((   node_memory_MemTotal_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_MemFree_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_Buffers_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} -   node_memory_Cached_bytes{job=\"node-exporter\", instance=\"x.x.x.x:xxxx\"} )) "
	LABEL_MASTER_ROLE                       = "node-role.kubernetes.io/master"
	URL_INSTALL_REPO                        = "https://github.com/tmax-cloud/install-helm-operator.git -b 5.0"
	HYPERCLOUD_CONSOLE_LABEL_APP            = "app"
	HYPERCLOUD_CONSOLE_LABEL_APP_KEY        = "console"
	HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD     = "hypercloud"
	HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD_KEY = "ui"
)

// func InstallPrometheus(res http.ResponseWriter, req *http.Request) {

// 	InstallPrometheusCommand()

// 	klog.Infoln("Success to exec install prometheus")
// 	util.SetResponse(res, "Success to exec install prometheus", nil, http.StatusInternalServerError)
// 	return
// }

func InstallCommand() {
	exec.Command("git", "clone", URL_INSTALL_REPO, "/installer/helm").Output()
	exec.Command("bash", "kubectl", "create", "namespace", "helm-ns").Output()
	exec.Command("bash", "kubectl", "apply", "-f", "/installer/helm/manifest").Output()
}

func HealthCheck() (*http.Response, error) {
	url := "http://prometheus-k8s.hypercloud5-system.svc.cluster.local:9090/-/healthy"
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // ignore certificate

	client := http.Client{
		Timeout: 15 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		klog.Errorln(err)
		return nil, err
	} else {
		return response, nil
	}
}
