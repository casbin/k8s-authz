package main

import (
	"testing"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateUnsupportedOperation(t *testing.T) {
	var gs CasbinServerHandler

	admissionReview := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{},
			Allowed: false,
		},
	}

	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Operation: "NONEXISTENT",
		},
	}

	gs.validate(&ar, admissionReview.Response)

	if admissionReview.Response.Result.Message != "Operation not supported" || admissionReview.Response.Allowed {
		t.Errorf("Invalid operation NONEXISTENT not rejected")
	}
}

func TestValidateCreateOperation(t *testing.T) {
	var gs CasbinServerHandler

	admissionReview := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{},
			Allowed: false,
		},
	}

	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Operation: "CREATE",
			Kind: metav1.GroupVersionKind{
				Kind: "NONEXISTENT",
			},
		},
	}

	gs.validate(&ar, admissionReview.Response)

	if admissionReview.Response.Result.Message == "Operation not supported" {
		t.Errorf("CREATE operation rejected")
	}
}

func TestValidateUpdateOperation(t *testing.T) {
	var gs CasbinServerHandler

	admissionReview := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{},
			Allowed: false,
		},
	}

	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Operation: "UPDATE",
			Kind: metav1.GroupVersionKind{
				Kind: "NONEXISTENT",
			},
		},
	}

	gs.validate(&ar, admissionReview.Response)

	if admissionReview.Response.Result.Message == "Operation not supported" {
		t.Errorf("UPDATE operation rejected")
	}
}

func TestValidateUnsupportedKind(t *testing.T) {
	var gs CasbinServerHandler

	admissionReview := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{},
			Allowed: false,
		},
	}

	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Operation: "CREATE",
			Kind: metav1.GroupVersionKind{
				Kind: "NONEXISTENT",
			},
		},
	}

	gs.validate(&ar, admissionReview.Response)

	if admissionReview.Response.Result.Message != "Kind not supported" || admissionReview.Response.Allowed {
		t.Errorf("Invalid kind NONEXISTENT not rejected")
	}
}

func TestValidatePodSecurityPolicy(t *testing.T) {
	var gs CasbinServerHandler

	admissionReview := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Result:  &metav1.Status{},
			Allowed: false,
		},
	}

	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Operation: "CREATE",
			Kind: metav1.GroupVersionKind{
				Kind: "PodSecurityPolicy",
			},
		},
	}

	gs.validate(&ar, admissionReview.Response)

	if admissionReview.Response.Result.Message == "Kind not supported" || admissionReview.Response.Allowed {
		t.Errorf("Valid kind PodSecurityPolicy rejected")
	}
}
