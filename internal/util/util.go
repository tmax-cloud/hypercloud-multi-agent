package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"regexp"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

type ClusterMemberInfo struct {
	Id          int64
	Namespace   string
	Cluster     string
	MemberId    string
	MemberName  string
	Attribute   string
	Role        string
	Status      string
	CreatedTime time.Time
	UpdatedTime time.Time
}

var (
	SMTPUsernamePath       string
	SMTPPasswordPath       string
	SMTPHost               string
	SMTPPort               int
	AccessSecretPath       string
	accessSecret           string
	username               string
	password               string
	inviteMail             string
	HtmlHomePath           string
	TokenExpiredDate       string
	ParsedTokenExpiredDate time.Duration
)

//Jsonpatch를 담을 수 있는 구조체
type PatchOps struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func ReadFile() {
	content, err := ioutil.ReadFile(AccessSecretPath)
	if err != nil {
		klog.Errorln(err)
		return
	}
	accessSecret = string(content)
	// klog.Infoln(accessSecret)

	content, err = ioutil.ReadFile(SMTPUsernamePath)
	if err != nil {
		klog.Errorln(err)
		return
	}
	username = string(content)

	content, err = ioutil.ReadFile(SMTPPasswordPath)
	if err != nil {
		klog.Errorln(err)
		return
	}
	password = string(content)

	ParsedTokenExpiredDate = parseDate(TokenExpiredDate)

}
func parseDate(tokenExpiredDate string) time.Duration {
	regex := regexp.MustCompile("[0-9]+")
	num := regex.FindAllString(tokenExpiredDate, -1)[0]
	parsedNum, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	regex = regexp.MustCompile("[a-z]+")
	unit := regex.FindAllString(tokenExpiredDate, -1)[0]

	switch unit {
	case "minutes":
		return time.Minute * time.Duration(parsedNum)
	case "hours":
		return time.Hour * time.Duration(parsedNum)
	case "days":
		return time.Hour * time.Duration(24) * time.Duration(parsedNum)
	case "weeks":
		return time.Hour * time.Duration(24) * time.Duration(7) * time.Duration(parsedNum)
	default:
		return time.Hour * time.Duration(24) * time.Duration(7) //1days
	}
}

// Jsonpatch를 하나 만들어서 slice에 추가하는 함수
func CreatePatch(po *[]PatchOps, o, p string, v interface{}) {
	*po = append(*po, PatchOps{
		Op:    o,
		Path:  p,
		Value: v,
	})
}

// Response.result.message에 err 메시지 넣고 반환
func ToAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func SetResponse(res http.ResponseWriter, outString string, outJson interface{}, status int) http.ResponseWriter {

	//set Cors
	// res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Max-Age", "3628800")
	res.Header().Set("Access-Control-Expose-Headers", "Content-Type, X-Requested-With, Accept, Authorization, Referer, User-Agent")

	//set Out
	if outJson != nil {
		res.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(outJson)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		//set StatusCode
		res.WriteHeader(status)
		res.Write(js)
		return res

	} else {
		//set StatusCode
		res.WriteHeader(status)
		res.Write([]byte(outString))
		return res

	}
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func Remove(slice []string, item string) []string {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	// for _, item := range items {
	if _, ok := set[item]; ok {
		delete(set, item)
	}
	// }

	var newSlice []string
	for k, _ := range set {
		newSlice = append(newSlice, k)
	}
	return newSlice
}

// func Remove(slice []string, items []string) []string {
// 	set := make(map[string]struct{}, len(slice))
// 	for _, s := range slice {
// 		set[s] = struct{}{}
// 	}

// 	for _, item := range items {
// 		_, ok := set[item]
// 		if ok {
// 			delete(set, item)
// 		}
// 	}

// 	var newSlice []string
// 	for k, _ := range set {
// 		newSlice = append(newSlice, k)
// 	}
// 	return newSlice
// }
