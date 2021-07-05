package prometheus

import (
	"crypto/tls"
	"net/http"
	"os/exec"
	"time"

	"github.com/tmax-cloud/hypercloud-multi-agent/internal/util"

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
	URL_INSTALL_REPO                        = "https://github.com/tmax-cloud/hypercloud-multi-agent-install-repo.git"
	HYPERCLOUD_CONSOLE_LABEL_APP            = "app"
	HYPERCLOUD_CONSOLE_LABEL_APP_KEY        = "console"
	HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD     = "hypercloud"
	HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD_KEY = "ui"
)

func InstallPrometheus(res http.ResponseWriter, req *http.Request) {

	InstallCommand()

	// if err := clonePrometheus(); err != nil {
	// 	msg := "Failed to clone prometheus. " + err.Error()
	// 	klog.Errorln(msg)
	// 	util.SetResponse(res, msg, nil, http.StatusInternalServerError)
	// 	return
	// }

	// if msg, err := exec.Command("chmod", "+x", "/tmp/git/prometheus/install.sh").Output(); err != nil {
	// 	klog.Errorln("Failed to chmod install.sh. \n" + string(msg))
	// 	util.SetResponse(res, "Failed to chmod install.sh \n"+string(msg), nil, http.StatusInternalServerError)
	// 	return
	// }

	// if msg, err := exec.Command("bash", "/tmp/git/prometheus/install.sh").Output(); err != nil {
	// 	klog.Errorln("Failed to exec install.sh \n" + string(msg))
	// 	util.SetResponse(res, "Failed to exec install.sh \n"+string(msg), nil, http.StatusInternalServerError)
	// 	return
	// }
	// prometheusIngress

	klog.Infoln("Success to exec install prometheus")
	util.SetResponse(res, "Success to exec install prometheus", nil, http.StatusInternalServerError)
	return
}

func InstallCommand() {
	exec.Command("git", "clone", URL_INSTALL_REPO, "/installer").Output()
	exec.Command("chmod", "+x", "/installer/main.sh").Output()
	exec.Command("bash", "/installer/main.sh").Output()
}

// func UnInstallPrometheus(res http.ResponseWriter, req *http.Request) {
// 	if _, err := os.Stat("/tmp/git/prometheus"); os.IsNotExist(err) {
// 		klog.Errorln("Prometheus install directory is removed." + err.Error())
// 		util.SetResponse(res, "Failed to exec uninstall.sh "+err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	if _, err := exec.Command("chmod", "+x", "/tmp/git/prometheus/uninstall.sh").Output(); err != nil {
// 		klog.Errorln("Failed to chmod uninstall.sh " + err.Error())
// 		util.SetResponse(res, "Failed to chmod uninstall.sh "+err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	if _, err := exec.Command("bash", "/tmp/git/prometheus/uninstall.sh").Output(); err != nil {
// 		klog.Errorln("Failed to exec uninstallin.sh " + err.Error())
// 		util.SetResponse(res, "Failed to exec uninstall.sh "+err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	if _, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/git/prometheus").Output(); err != nil {
// 		klog.Errorln("Failed to remove prometheus install directory " + err.Error())
// 		util.SetResponse(res, "Failed to remove prometheus install directory "+err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}
// 	klog.Infoln("Success to exec uninstall prometheus")
// 	util.SetResponse(res, "Success to exec uninstall prometheus", nil, http.StatusInternalServerError)
// 	return
// }

// func deletePrometheusDir() error {
// 	if _, err := os.Stat("/tmp/git/prometheus"); os.IsNotExist(err) {
// 		return nil
// 	} else if err != nil {
// 		return err
// 	} else {
// 		if _, err := exec.Command("/bin/sh", "-c", "rm -rf /tmp/git/prometheus").Output(); err != nil {
// 			return err
// 		} else {
// 			return nil
// 		}
// 	}
// }

// func clonePrometheus() error {
// 	if _, err := os.Stat("/tmp/git/prometheus"); os.IsNotExist(err) {
// 		_, err := git.PlainClone("/tmp/git/prometheus", false, &git.CloneOptions{
// 			URL:             "https://github.com/tmax-cloud/install-prometheus.git",
// 			ReferenceName:   plumbing.ReferenceName("refs/heads/5.0-agent"),
// 			InsecureSkipTLS: true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		return errors.New("Prometheus git directory is already existed")
// 	}
// 	return nil
// }

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
