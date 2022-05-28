package handler

import (
	"fmt"
	"io/ioutil"

	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/casbin/k8s-authz/internal/model"
	"github.com/gin-gonic/gin"
)

//Main Handler
func Handler(c *gin.Context) {

	data, _ := ioutil.ReadAll(c.Request.Body)
	var admissionReview admission.AdmissionReview
	var decoder runtime.Decoder = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	decoder.Decode(data, nil, &admissionReview)

	//for debug only. Todo:remove this block of code
	if admissionReview.Request.Namespace != "default" {
		approveResponse(c, string(admissionReview.Request.UID))
		return
	}
	//fmt.Println(string(data))
	fmt.Printf("%s\n", admissionReview.Request.Resource.String())

	//currently we are going to handle these resources:
	uid := admissionReview.Request.UID
	resource := admissionReview.Request.Resource.Resource

	switch resource {
	case "deployments":
		model.MountDeploymentObject(&admissionReview)
	}

	err := model.EnforcerList.Enforce(&admissionReview)
	if err != nil {
		fmt.Println("rejected")
		rejectResponse(c, string(uid), err.Error())
		return
	}

	fmt.Println("approved")
	approveResponse(c, string(uid))

}

func rejectResponse(c *gin.Context, uid string, rejectReason string) {
	c.JSON(200, gin.H{
		"apiVersion": "admission.k8s.io/v1",
		"kind":       "AdmissionReview",
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": false,
			"status": map[string]interface{}{
				"code":    403,
				"message": rejectReason,
			},
		},
	})
}

func approveResponse(c *gin.Context, uid string) {
	c.JSON(200, gin.H{
		"apiVersion": "admission.k8s.io/v1",
		"kind":       "AdmissionReview",
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": true,
		},
	})
}
