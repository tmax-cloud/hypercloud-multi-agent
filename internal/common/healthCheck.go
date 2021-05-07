package common

import (
	"io/ioutil"
	"net/http"

	"github.com/tmax-cloud/hypercloud-multi-agent/internal/prometheus"
	"github.com/tmax-cloud/hypercloud-multi-agent/internal/util"
	"k8s.io/klog"
)

type health struct {
	ModuleName string `json:"moduleName"`
	Ready      bool   `json:"ready"`
	Status     int    `json:"status"`
	Msg        string `json:"msg"`
}

func HealthCheck(res http.ResponseWriter, req *http.Request) {
	ret := []health{}

	if response, err := prometheus.HealthCheck(); err != nil {
		klog.Errorln("Failed to get prometheus status " + err.Error())
		util.SetResponse(res, "Failed to get prometheus status "+err.Error(), nil, http.StatusInternalServerError)
		return
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		prometheusHealth := health{}
		prometheusHealth.ModuleName = "prometheus"
		prometheusHealth.Msg = string(data)
		prometheusHealth.Status = response.StatusCode
		if response.StatusCode >= 200 && response.StatusCode < 400 {
			prometheusHealth.Ready = true
		}
		ret = append(ret, prometheusHealth)
	}
	klog.Infoln("Success to get module status")
	util.SetResponse(res, "Success to get module status", ret, http.StatusOK)
	return
}
