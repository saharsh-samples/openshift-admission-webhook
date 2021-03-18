package model

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestAdmissionReview_Validate(t *testing.T) {
	test.AssertTrue("Expected nil error", (&AdmissionReview{}).Validate() == nil, t)
}
