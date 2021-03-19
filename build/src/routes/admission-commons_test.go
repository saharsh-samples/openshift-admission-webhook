package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

func Test_allowWithPatches_for_multiple_patches(t *testing.T) {

	// Arrange
	mockW := httptest.NewRecorder()
	helper := &admissionHelper{
		jsonUtils: utils.Bootstrap(nil).JSONUtils,
		w:         mockW,
	}

	// Act
	helper.allowWithPatches(types.UID("some-uid"), []string{"patch-1", "patch-2"})

	// Assert
	resp := mockW.Result()
	test.AssertEquals("", http.StatusOK, resp.StatusCode, t)

	body, _ := ioutil.ReadAll(resp.Body)
	response := &v1beta1.AdmissionReview{}
	err := json.Unmarshal(body, response)
	test.AssertTrue("expected no json error with response", err == nil, t)
	test.AssertFalse("Expected response to be non-nil", response.Response == nil, t)
	test.AssertEquals("", types.UID("some-uid"), response.Response.UID, t)
	test.AssertTrue("", response.Response.Allowed, t)
	test.AssertTrue("Expected Result to be nil", response.Response.Result == nil, t)
	test.AssertEquals("", v1beta1.PatchTypeJSONPatch, *response.Response.PatchType, t)
	test.AssertEquals("", `[ patch-1, patch-2 ]`, string(response.Response.Patch), t)

}
