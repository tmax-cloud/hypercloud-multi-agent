package admission

import (
	"errors"

	// jsonpatch "github.com/evanphx/json-patch"

	"k8s.io/api/admission/v1beta1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	KUBE_SYSTEM_NAMESPACE string = "kube-system"
	// PROMETHEUS_SYSTEM_NAMESPACE string = "monitoring"
	HYPERCLOUD_SYSTEM_NAMESPACE string = "hypercloud5-system"
	// HYPERCLOUD_SYSTEM_ADMIN     string = "hypercloud5-admin"
	HYPERCLOUD_SYSTEM_ADMIN  string = "kubernetes-admin"
	HYPERCLOUD_DEFAULT_GROUP string = "hypercloud5"
)

func DenyRequest(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	// reviewResponse := v1beta1.AdmissionResponse{}
	userInfo := ar.Request.UserInfo
	// requestNamespace := ar.Request.Namespace

	if userInfo.Username == HYPERCLOUD_SYSTEM_ADMIN {
		return ToAdmissionResponse(nil)
	}

	// if util.Contains(userInfo.Groups, HYPERCLOUD_DEFAULT_GROUP) {
	msg := "Can not " + string(ar.Request.Operation) + " requests to the " + KUBE_SYSTEM_NAMESPACE + "or " + HYPERCLOUD_SYSTEM_NAMESPACE + "namespaces"
	klog.Info(msg)
	return ToAdmissionResponse(errors.New(msg))
	// }

	// if requestNamespace == KUBE_SYSTEM_NAMESPACE || requestNamespace == PROMETHEUS_SYSTEM_NAMESPACE || requestNamespace == HYPERCLOUD_SYSTEM_NAMESPACE {
	// 	return ToAdmissionResponse(errors.New("Can not send any requests to the "+ KUBE_SYSTEM_NAMESPACE + "or " + PROMETHEUS_SYSTEM_NAMESPACE + "namespaces"))
	// }

	// reviewResponse.Allowed = true

	// return &reviewResponse
}
