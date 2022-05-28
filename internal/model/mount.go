package model

import (
	"encoding/json"

	admission "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
)

func MountDeploymentObject(admissionReview *admission.AdmissionReview) error {
	admissionReview.Request.Object.Object = nil
	if len(admissionReview.Request.Object.Raw) != 0 {
		var deploymentObject app.Deployment
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &deploymentObject)
		if err != nil {
			return err
		}
		admissionReview.Request.Object.Object = &deploymentObject
	}

	admissionReview.Request.OldObject.Object = nil
	if len(admissionReview.Request.OldObject.Raw) != 0 {
		var deploymentOldObject app.Deployment
		err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &deploymentOldObject)
		if err != nil {
			return err
		}
		admissionReview.Request.OldObject.Object = &deploymentOldObject
	}
	return nil
}
