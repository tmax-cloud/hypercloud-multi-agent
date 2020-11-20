/*


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

package controllers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/tmax-cloud/hypercloud-multi-agent/util"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// ServiceReconciler reconciles a Memcached object
type ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=service;,verbs=get;update;patch;list;watch;create;post;delete;

func (r *ServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()

	//get service
	svc := &corev1.Service{}
	if err := r.Get(context.TODO(), req.NamespacedName, svc); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Service resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}

		r.Log.Error(err, "Failed to get service")
		return ctrl.Result{}, err
	}

	// wait for configmap setting
	if strings.Compare(util.MGMTIP, "") == 0 {
		r.Log.Info("Wait for configmap setting ...")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Duration(10 * time.Second)}, nil
	}

	// delete handling
	if !svc.DeletionTimestamp.IsZero() {
		_ = putRequestURL("")
	}

	// send url to MGMT
	if err := putRequestURL(getURL(*svc)); err != nil {
		r.Log.Error(err, "putRequestURL error")
	}

	return ctrl.Result{}, nil
}

func getURL(svc corev1.Service) string {
	hyperURL := ""

	if svc.Spec.Type == "NodePort" {
		for _, node := range util.MYNODEINFO {
			hyperURL = hyperURL + "https://" + node.Ip + ":" + strconv.Itoa(int(svc.Spec.Ports[0].NodePort)) + ";"
		}
	}

	if svc.Spec.Type == "LoadBalancer" {
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			hyperURL = "https://" + ingress.IP + ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port)) + ";"
		}
	}

	return hyperURL
}

func putRequestURL(hypercloudurl string) error {
	req, _ := http.NewRequest("PUT", util.URL_PREFIX+util.MGMTIP+":"+util.MGMTPORT+util.URL_HTPERCLOUD_URL_PATH, nil)
	qeuryParam := url.Values{}
	qeuryParam.Add(util.CONFIGMAP_MY_CLUSTERNAME, util.MYCLUSTERNAME)
	qeuryParam.Add(util.URL_PARAM_URL, hypercloudurl)
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

func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	_, err := ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		WithEventFilter(
			predicate.Funcs{
				CreateFunc: func(e event.CreateEvent) bool {
					svc := e.Object.(*corev1.Service).DeepCopy()
					return checkValidSelector(*svc)
				},
				UpdateFunc: func(e event.UpdateEvent) bool {
					svc := e.ObjectNew.(*corev1.Service).DeepCopy()
					return checkValidSelector(*svc) && checkValidUpdate(e)
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					return false
				},
				GenericFunc: func(e event.GenericEvent) bool {
					return false
				},
			},
		).
		Build(r)

	return err
}

func checkValidUpdate(e event.UpdateEvent) bool {
	oldsvc := e.ObjectOld.(*corev1.Service).DeepCopy()
	newsvc := e.ObjectNew.(*corev1.Service).DeepCopy()

	if oldsvc.Spec.Ports[0].NodePort != newsvc.Spec.Ports[0].NodePort ||
		strings.Compare(string(oldsvc.Spec.Type), string(newsvc.Spec.Type)) != 0 ||
		len(oldsvc.Status.LoadBalancer.Ingress) != len(newsvc.Status.LoadBalancer.Ingress) ||
		strings.Compare(oldsvc.Status.LoadBalancer.Ingress[0].IP, newsvc.Status.LoadBalancer.Ingress[0].IP) != 0 {
		return true
	}

	return false
}

func checkValidSelector(svc corev1.Service) bool {
	selector := svc.Spec.Selector
	if selector != nil {
		if val, ok := selector[util.HYPERCLOUD_CONSOLE_LABEL_APP]; ok && strings.Compare(val, util.HYPERCLOUD_CONSOLE_LABEL_APP_KEY) == 0 {
			if val, ok := selector[util.HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD]; ok && strings.Compare(val, util.HYPERCLOUD_CONSOLE_LABEL_HYPERCLOUD_KEY) == 0 {
				return true
			}
		}
	}

	return false
}
