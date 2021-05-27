package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/casbin/casbin/v2"
)

type CasbinServerHandler struct {
}

func (gs *CasbinServerHandler) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	glog.Info("Received request")

	if r.URL.Path != "/validate" {
		glog.Error("no validate")
		http.Error(w, "no validate", http.StatusBadRequest)
		return
	}

	arRequest := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &arRequest); err != nil {
		glog.Error("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	raw := arRequest.Request.Object.Raw
	user := arRequest.userInfo.username
	operation_name := arRequest.operation

	if err := json.Unmarshal(raw, &user); err != nil {
		glog.Error("error deserializing User name")
		return
	}
	if err := json.Unmarshal(raw, &operation_name); err != nil {
		glog.Error("error deserializing Operation name")
		return
	}
	
	e, err := casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
	if err != nil {
		glog.Errorf("Filed to load the policies: %v", err)
		return
	}

	if e.HasPermissionForUser(user, []string{operation_name}) = true {
	response := v1.AdmissionReview{
		Response: &v1.AdmissionResponse{
			Allowed: true,
		},
		}
	}
	response := v1.AdmissionReview{
		Response: &v1.AdmissionResponse{
		  Allowed: false,
		  Result: &metav1.Status{
			Message: " You are not authorized to perform any operations on these pods!",
		  },
		},
	  }

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
