package model

import (
	"fmt"

	"k8s.io/api/admission/v1beta1"
)

// AdmissionReview makes v1beta1.AdmissionReview compatible with
// go-mux-sql-starter's JSONBody interface.
type AdmissionReview struct {
	*v1beta1.AdmissionReview
}

// Validate for instances of AdmissionReview
func (body *AdmissionReview) Validate() error {

	errors := make([]interface{}, 0)

	if body.AdmissionReview == nil || body.Request == nil {
		errors = append(errors, "Request must be non-nil")
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}
	return nil

}
