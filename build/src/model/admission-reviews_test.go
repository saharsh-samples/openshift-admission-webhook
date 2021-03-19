package model

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
	"k8s.io/api/admission/v1beta1"
)

func TestAdmissionReview_Validate(t *testing.T) {

	body := &AdmissionReview{}

	test.AssertFalse("Expected validation error", body.Validate() == nil, t)
	test.AssertEquals("", "[Request must be non-nil]", body.Validate().Error(), t)

	body.AdmissionReview = &v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{}}
	test.AssertTrue("Expected nil error", body.Validate() == nil, t)
}
