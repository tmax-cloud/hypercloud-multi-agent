package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hyperv1 "github.com/tmax-cloud/hypercloud-multi-agent/external/hyper/v1"
	prometheus "github.com/tmax-cloud/hypercloud-multi-agent/external/prometheus"
)

// +kubebuilder:rbac:groups="",resources=nodes;configmaps;services,verbs=get;list;watch;create;update;patch;delete

var (
	MGMTIP               string
	MGMTPORT             string
	MYCLUSTERNAME        string
	prometheusPort       string
	prometheusExportPort string
	myNodes              = &corev1.NodeList{}
	MYNODEINFO           []hyperv1.NodeInfo
)

type Collector struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (c *Collector) Collect() {
	c.setConfigEnv()

	//update hcr's resource usage every 10 seconds
	go func() {
		for {
			c.handleHCR()
			time.Sleep(10 * time.Second)
		}
	}()

	//update node info every 10 seconds
	go func() {
		for {
			time.Sleep(60 * time.Second)
			c.updateMyNodeInfo()
		}
	}()
}

func (c *Collector) handleHCR() {
	c.updateHCR()

	if err := putRequestHCR(MYNODEINFO); err != nil {
		klog.Error(err)
	}

	klog.Info(MYNODEINFO)
}

func (c *Collector) updateHCR() {
	cleanNodeInfoResource()

	requestPrometheusPod()
	requestPrometheusCPU()
	requestPrometheusStorage()
	requestPrometheusMemory()
}

func requestPrometheusMemory() {
	for index, node := range myNodes.Items {
		resourceMemory := hyperv1.ResourceType{
			Type:     "memory",
			Capacity: "",
			Usage:    "",
		}

		memoryCap, _ := node.Status.Capacity.Memory().AsInt64()
		resourceMemory.Capacity = strconv.FormatInt(memoryCap, 10)
		resourceMemory.Usage = requestPrometheusQuery(PROMETHEUS_QUERY_MEMORY_USAGE, MYNODEINFO[index].Ip)

		MYNODEINFO[index].Resources = append(MYNODEINFO[index].Resources, resourceMemory)
	}
}

func requestPrometheusStorage() {
	for index, node := range myNodes.Items {
		resourceStorage := hyperv1.ResourceType{
			Type:     "storage",
			Capacity: "",
			Usage:    "",
		}
		StorageCap, _ := node.Status.Capacity.StorageEphemeral().AsInt64()
		resourceStorage.Capacity = strconv.FormatInt(StorageCap, 10)
		resourceStorage.Usage = requestPrometheusQuery(PROMETHEUS_QUERY_STORAGE_USAGE, MYNODEINFO[index].Ip)

		MYNODEINFO[index].Resources = append(MYNODEINFO[index].Resources, resourceStorage)
	}
}

func requestPrometheusCPU() {
	for index, node := range myNodes.Items {
		resourceCPU := hyperv1.ResourceType{
			Type:     "cpu",
			Capacity: "",
			Usage:    "",
		}
		CPUCap, _ := node.Status.Capacity.Cpu().AsInt64()
		resourceCPU.Capacity = strconv.FormatInt(CPUCap, 10)
		resourceCPU.Usage = requestPrometheusQuery(PROMETHEUS_QUERY_CPU_USAGE, MYNODEINFO[index].Ip)

		MYNODEINFO[index].Resources = append(MYNODEINFO[index].Resources, resourceCPU)
	}
}

func requestPrometheusPod() {
	for index, node := range myNodes.Items {
		resourcePod := hyperv1.ResourceType{
			Type:     "pod",
			Capacity: "",
			Usage:    "",
		}
		podCap, _ := node.Status.Capacity.Pods().AsInt64()
		resourcePod.Capacity = strconv.FormatInt(podCap, 10)
		resourcePod.Usage = requestPrometheusQuery(PROMETHEUS_QUERY_POD_USAGE, MYNODEINFO[index].Ip)

		MYNODEINFO[index].Resources = append(MYNODEINFO[index].Resources, resourcePod)
	}
}

func requestPrometheusQuery(query string, nodeIp string) string {
	usage := 0.0
	req, err := http.NewRequest("GET", URL_PREFIX+myNodes.Items[0].Status.Addresses[0].Address+":"+prometheusPort+PROMETHEUS_QUERY_PATH, nil)
	if err != nil {
		klog.Info(err)
	}

	qeuryParam := url.Values{}
	qeuryParam.Add(PROMETHEUS_QUERY_KEY_QUERY, strings.Replace(strings.Replace(query, "x.x.x.x:xxxx", nodeIp+":"+prometheusExportPort, -1), "y.y.y.y", nodeIp, -1))
	qeuryParam.Add(PROMETHEUS_QUERY_KEY_TIME, strconv.FormatInt(time.Now().Unix(), 10))

	req.URL.RawQuery = qeuryParam.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	payloadbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	res := &prometheus.Responese{}
	json.Unmarshal(payloadbytes, res)
	for _, value := range res.Data.Result {
		now, _ := strconv.ParseFloat(value.Value[1], 64)
		usage += now
	}

	usageStr := fmt.Sprintf("%f", usage)

	return usageStr
}

func putRequestHCR(nodeInfo []hyperv1.NodeInfo) error {
	payloadbytes, _ := json.Marshal(nodeInfo)
	req, _ := http.NewRequest("PUT", URL_PREFIX+MGMTIP+":"+MGMTPORT+URL_HYPERCLUSTERRESOURCE_PATH, bytes.NewBuffer(payloadbytes))
	qeuryParam := url.Values{}
	qeuryParam.Add(CONFIGMAP_MY_CLUSTERNAME, MYCLUSTERNAME)
	req.URL.RawQuery = qeuryParam.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Collector) setConfigEnv() {
	for {
		cm := &corev1.ConfigMap{}

		if err := c.Get(context.TODO(), types.NamespacedName{Name: CONFIGMAP_NAME, Namespace: CONFIGMAP_NAMESPACE}, cm); err != nil {
			klog.Info("configmap doesn't exist in " + CONFIGMAP_NAMESPACE + " with " + CONFIGMAP_NAME + " name")
			time.Sleep(time.Second * 30)
			continue
		}

		MGMTIP = cm.Data[CONFIGMAP_MGNT_IP]
		MGMTPORT = cm.Data[CONFIGMAP_MGNT_PORT]
		MYCLUSTERNAME = cm.Data[CONFIGMAP_MY_CLUSTERNAME]

		break
	}

	for {
		svc := &corev1.Service{}

		if err := c.Get(context.TODO(), types.NamespacedName{Name: PROMETHEUS_SERVICE_NAME, Namespace: PROMETHEUS_SERVICE_NAMESPACE}, svc); err != nil {
			klog.Info("there is no prometheus service in " + PROMETHEUS_SERVICE_NAMESPACE + " with " + PROMETHEUS_SERVICE_NAME + " name")
			time.Sleep(time.Second * 30)
			continue
		}

		for _, port := range svc.Spec.Ports {
			if port.Name == "web" {
				prometheusPort = strconv.Itoa(int(port.NodePort))
			}
		}

		if err := c.Get(context.TODO(), types.NamespacedName{Name: PROMETHEUS_NODE_EXPORT_SERVICE_NAME, Namespace: PROMETHEUS_SERVICE_NAMESPACE}, svc); err != nil {
			klog.Info("there is no prometheus node-exporter in " + PROMETHEUS_SERVICE_NAMESPACE + " with " + PROMETHEUS_NODE_EXPORT_SERVICE_NAME + " name")
			time.Sleep(time.Second * 30)
			continue
		}

		for _, port := range svc.Spec.Ports {
			if port.Name == "https" {
				prometheusExportPort = strconv.Itoa(int(port.Port))
			}
		}

		break
	}

	c.updateMyNodeInfo()
}

func (c *Collector) updateMyNodeInfo() {
	c.List(context.TODO(), myNodes)
	MYNODEINFO = []hyperv1.NodeInfo{}
	for _, node := range myNodes.Items {
		nodeInfo := hyperv1.NodeInfo{
			Name:      node.Name,
			Ip:        node.Status.Addresses[0].Address,
			IsMaster:  checkMaster(node),
			Resources: []hyperv1.ResourceType{},
		}
		MYNODEINFO = append(MYNODEINFO, nodeInfo)
	}
}

func checkMaster(node corev1.Node) bool {
	if _, has := node.Labels[LABEL_MASTER_ROLE]; has {
		return true
	}
	return false
}

func cleanNodeInfoResource() {
	for index, _ := range MYNODEINFO {
		MYNODEINFO[index].Resources = []hyperv1.ResourceType{}
	}
}
