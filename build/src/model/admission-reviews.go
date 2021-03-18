package model

import (
	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
	"k8s.io/api/admission/v1beta1"
)

// AdmissionReview makes v1beta1.AdmissionReview compatible with
// go-mux-sql-starter's JSONBody interface.
type AdmissionReview struct {
	*v1beta1.AdmissionReview
	*utils.AlwaysValidJSON
}
