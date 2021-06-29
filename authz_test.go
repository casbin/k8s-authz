package main

import (
	"testing"
	//"encoding/json"

//	"k8s.io/api/admission/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/runtime"
//	"github.com/casbin/casbin/v2"
	//	authenticationv1 "k8s.io/api/authentication/v1"
	// "github.com/casbin/k8s-authz/server"
	// "github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)


//var test_operation = "CREATE"

/*var (
	AdmissionRequestpod = v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &v1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000001",
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: []byte(`{"metadata": {
        						"name": "test",
        						"uid": "e911857d-c318-11e8-bbad-025000000001",
						        "creationTimestamp": "2018-09-28T12:20:39Z"
      						}}`),
			},
			UserInfo: {

				"username": "test_user",

			},
		},
	}
)
*/

/*
func TestPolicy(t *testing.T) {
	//declare variable name of user and operation and then enforce it with casbin to check the policy verifications
	e, err := casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
	e.AddPermissionForUser("test_user", "CREATE")
	cs := CasbinServerHandler{}

	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "c7m",
					"namespace": "yolo",
					"creationTimestamp": null,
					"labels": {
						"name": "c7m"
					},
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [
						{
							"name": "default-token-5z7xl",
							"secret": {
								"secretName": "default-token-5z7xl"
							}
						}
					],
					"containers": [
						{
							"name": "c7m",
							"image": "centos:7",
							"command": [
								"/bin/bash"
							],
							"args": [
								"-c",
								"trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"
							],
							"resources": {},
							"volumeMounts": [
								{
									"name": "default-token-5z7xl",
									"readOnly": true,
									"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
								}
							],
							"terminationMessagePath": "/dev/termination-log",
							"terminationMessagePolicy": "File",
							"imagePullPolicy": "IfNotPresent"
						}
					],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [
						{
							"key": "node.kubernetes.io/not-ready",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						},
						{
							"key": "node.kubernetes.io/unreachable",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						}
					],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	ar := v1.AdmissionReview{
		Request: &v1.AdmissionRequest{
			Operation: "NONEXISTENT",
		},
	}
	resp := v1.AdmissionRequest{}

	//	r := v1beta1.AdmissionReview{}
	//	w := rawJSON
	//	response := cs.serve(w,r)
	resp.UserInfo.Username = user
	handler := http.HandlerFunc(cs.serve)
	handler.ServeHTTP([]byte(resp), ar)
}
*/



func TestValidationHandler(t *testing.T) {
	cs := CasbinServerHandler{}
	r, err := http.NewRequest("GET", "/validate", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(cs.serve)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}
}
