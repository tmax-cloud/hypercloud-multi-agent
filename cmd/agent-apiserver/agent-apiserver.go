package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	gmux "github.com/gorilla/mux"
	// "github.com/tmax-cloud/hypercloud-api-server/cluster"

	"github.com/tmax-cloud/hypercloud-multi-agent/internal/common"

	// "github.com/tmax-cloud/hypercloud-multi-agent/cluster"

	"github.com/tmax-cloud/hypercloud-multi-agent/internal/prometheus"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"

	"net/http"

	"github.com/robfig/cron"
	//kafkaConsumer "github.com/tmax-cloud/hypercloud-multi-agent/util/consumer"
)

type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

var (
	port int

// 	certFile string
// 	keyFile  string
)

const (
	PROMETHEUS_MODULE_NAME string = "prometheus"
)

func main() {
	// For tls
	flag.IntVar(&port, "port", 8080, "hypercloud5-multi-agent port")
	// flag.StringVar(&certFile, "certFile", "/run/secrets/tls/hypercloud-multi-agent.crt", "hypercloud5-multi-agent cert")
	// flag.StringVar(&keyFile, "keyFile", "/run/secrets/tls/hypercloud-multi-agent.key", "hypercloud5-multi-agent key")
	// flag.StringVar(&admission.SidecarContainerImage, "sidecarImage", "fluent/fluent-bit:1.5-debug", "Fluent-bit image name.")
	// flag.StringVar(&util.SMTPHost, "smtpHost", "mail.tmax.co.kr", "SMTP Server Host Address")
	// flag.IntVar(&util.SMTPPort, "smtpPort", 25, "SMTP Server Port")
	// flag.StringVar(&util.SMTPUsernamePath, "smtpUsername", "/run/secrets/smtp/username", "SMTP Server Username")
	// flag.StringVar(&util.SMTPPasswordPath, "smtpPassword", "/run/secrets/smtp/password", "SMTP Server Password")
	// flag.StringVar(&util.AccessSecretPath, "accessSecret", "/run/secrets/token/accessSecret", "Token Access Secret")
	// flag.StringVar(&util.HtmlHomePath, "htmlPath", "/run/configs/html/", "Invite htlm path")
	// flag.StringVar(&util.TokenExpiredDate, "tokenExpiredDate", "24hours", "Token Expired Date")

	// go util.ReadFile()

	// Get Hypercloud Operating Mode!!!
	// hcMode := os.Getenv("HC_MODE")
	// util.TokenExpiredDate = os.Getenv("INVITATION_TOKEN_EXPIRED_DATE")

	// For Log file
	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()

	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", os.ModeDir)
	}

	file, err := os.OpenFile(
		"./logs/multi-agent.log",
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644),
	)
	if err != nil {
		klog.Error(err, "Error Open", "./logs/multi-agent")
		return
	}
	defer file.Close()
	w := io.MultiWriter(file, os.Stdout)
	klog.SetOutput(w)

	// Logging Cron Job
	cronJob := cron.New()
	cronJob.AddFunc("1 0 0 * * ?", func() {
		input, err := ioutil.ReadFile("./logs/multi-agent.log")
		if err != nil {
			klog.Error(err)
			return
		}
		err = ioutil.WriteFile("./logs/multi-agent"+time.Now().Format("2006-01-02")+".log", input, 0644)
		if err != nil {
			klog.Error(err, "Error creating", "./logs/multi-agent")
			return
		}
		klog.Info("Log BackUp Success")
		os.Truncate("./logs/multi-agent.log", 0)
		file.Seek(0, os.SEEK_SET)
	})

	prometheus.InstallCommand()

	// keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
	// if err != nil {
	// 	klog.Errorf("Failed to load key pair: %s", err)
	// }

	// Req multiplexer
	mux := gmux.NewRouter()

	mux.HandleFunc("/install/{module}", serveInstall)
	// mux.HandleFunc("/uninstall/{module}", serveUninstall)
	// mux.HandleFunc("/update/{module}", serveUpdate)
	// mux.HandleFunc("/validation", serveValidation)
	mux.HandleFunc("/healthy", serveModuleHealthy)
	mux.HandleFunc("/livez", serveLivez)

	// HTTP Server Start
	klog.Info("Starting Hypercloud5-Agent server...")
	klog.Flush()

	whsvr := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		// TLSConfig: &tls.Config{Certificates: []tls.Certificate{keyPair}},
	}

	if err := whsvr.ListenAndServe(); err != nil {
		klog.Errorf("Failed to listen and serve Hypercloud5-Agent server: %s", err)
	}
}

func serveLivez(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
	}
}

func serveModuleHealthy(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		common.HealthCheck(res, req)
		break
	}
}

// func serveValidation(w http.ResponseWriter, r *http.Request) {
// 	klog.Infof("Http request: method=%s, uri=%s", r.Method, r.URL.Path)
// 	serve(w, r, admission.DenyRequest)
// }

func serveInstall(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	vars := gmux.Vars(req)
	switch req.Method {
	// case http.MethodGet:
	// cluster.ListInvitation(res, req)
	// break
	case http.MethodPost:
		if vars["module"] == PROMETHEUS_MODULE_NAME {
			prometheus.InstallPrometheus(res, req)
		} else {
			// errror
		}
		break
	default:
	}
}

// func serveUninstall(res http.ResponseWriter, req *http.Request) {
// 	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
// 	vars := gmux.Vars(req)
// 	switch req.Method {
// 	case http.MethodGet:
// 		// cluster.ListInvitation(res, req)
// 		break
// 	case http.MethodPost:
// 		if vars["module"] == PROMETHEUS_MODULE_NAME {
// 			prometheus.UnInstallPrometheus(res, req)
// 		} else {
// 			// errror
// 		}
// 		break
// 	default:
// 	}
// }

// func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
// 	var body []byte
// 	if r.Body != nil {
// 		if data, err := ioutil.ReadAll(r.Body); err == nil {
// 			body = data
// 		}
// 	}

// 	contentType := r.Header.Get("Content-Type")
// 	if contentType != "application/json" {
// 		klog.Errorf("contentType=%s, expect application/json", contentType)
// 		return
// 	}

// 	requestedAdmissionReview := v1beta1.AdmissionReview{}
// 	responseAdmissionReview := v1beta1.AdmissionReview{}

// 	if err := json.Unmarshal(body, &requestedAdmissionReview); err != nil {
// 		klog.Error(err)
// 		responseAdmissionReview.Response = admission.ToAdmissionResponse(err)
// 	} else {
// 		responseAdmissionReview.Response = admit(requestedAdmissionReview)
// 	}

// 	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

// 	respBytes, err := json.Marshal(responseAdmissionReview)

// 	klog.Infof("Response body: %s\n", respBytes)

// 	if err != nil {
// 		klog.Error(err)
// 		responseAdmissionReview.Response = admission.ToAdmissionResponse(err)
// 	}
// 	if _, err := w.Write(respBytes); err != nil {
// 		klog.Error(err)
// 		responseAdmissionReview.Response = admission.ToAdmissionResponse(err)
// 	}
// }
