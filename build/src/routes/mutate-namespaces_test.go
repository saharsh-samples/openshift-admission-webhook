package routes

import (
	"encoding/json"
	"fmt"
	"maw/integrations"
	"net/http"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	httpTest "github.com/saharsh-samples/go-mux-sql-starter/http/test"
	httpUtils "github.com/saharsh-samples/go-mux-sql-starter/http/utils"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestRegister_for_namespaceMutator(t *testing.T) {

	// Arrage
	agent := httpTest.NewMockRoutesAgent()
	resource := &namespaceMutator{}

	// Act
	resource.Register(agent)

	// Assert
	agent.VerifyThatRoute(t, "/admissions/namespaces").ForHTTPMethod(http.MethodPost).UsesHandler(resource.MutateNamespace)

}

func createNamespaceMutatorForTesting() *namespaceMutator {
	return &namespaceMutator{
		specialProvider: integrations.Bootstrap(nil).SpecialProvider,
		jsonUtils:       httpUtils.Bootstrap(nil).JSONUtils,
	}
}

func TestMutateNamespace_when_namespace_is_special(t *testing.T) {

	// Arrange
	mutator := createNamespaceMutatorForTesting()

	serverPort, appCtx := httpTest.StartTestServer([]base.Routes{mutator}, t)
	defer httpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/admissions/namespaces", serverPort)

	request := &v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			UID: types.UID("some-uid"),
			Object: runtime.RawExtension{
				Raw: json.RawMessage(`{"metadata":{"name":"special-namespace"}}`),
			},
		},
	}

	// Act
	response := &v1beta1.AdmissionReview{}
	status := httpTest.Post(t, url, request, response, "")

	// Assert
	test.AssertEquals("", http.StatusOK, status, t)

	test.AssertFalse("Expected response to be non-nil", response.Response == nil, t)
	test.AssertEquals("", request.Request.UID, response.Response.UID, t)
	test.AssertTrue("", response.Response.Allowed, t)
	test.AssertTrue("Expected Result to be nil", response.Response.Result == nil, t)
	test.AssertEquals("", v1beta1.PatchTypeJSONPatch, *response.Response.PatchType, t)
	test.AssertEquals("", `[ { "op" : "add", "path": "/metadata/annotations/openshift.io~1node-selector", "value": "workload-type=special"} ]`, string(response.Response.Patch), t)

}

func TestMutateNamespace_when_namespace_is_not_special(t *testing.T) {

	// Arrange
	mutator := createNamespaceMutatorForTesting()

	serverPort, appCtx := httpTest.StartTestServer([]base.Routes{mutator}, t)
	defer httpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/admissions/namespaces", serverPort)

	request := &v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			UID: types.UID("some-uid"),
			Object: runtime.RawExtension{
				Raw: json.RawMessage(`{"metadata":{"name":"ordinary-namespace"}}`),
			},
		},
	}

	// Act
	response := &v1beta1.AdmissionReview{}
	status := httpTest.Post(t, url, request, response, "")

	// Assert
	test.AssertEquals("", http.StatusOK, status, t)

	test.AssertFalse("Expected response to be non-nil", response.Response == nil, t)
	test.AssertEquals("", request.Request.UID, response.Response.UID, t)
	test.AssertTrue("", response.Response.Allowed, t)
	test.AssertTrue("Expected Result to be nil", response.Response.Result == nil, t)
	test.AssertTrue("Expected PatchType to be nil", response.Response.PatchType == nil, t)
	test.AssertTrue("Expected Patch to be nil", response.Response.Patch == nil, t)

}

func TestMutateNamespace_when_namespace_name_missing(t *testing.T) {

	// Arrange
	mutator := createNamespaceMutatorForTesting()

	serverPort, appCtx := httpTest.StartTestServer([]base.Routes{mutator}, t)
	defer httpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/admissions/namespaces", serverPort)

	request := &v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			UID: types.UID("some-uid"),
			Object: runtime.RawExtension{
				Raw: json.RawMessage(`{"metadata":{"key":"value"}}`),
			},
		},
	}

	// Act
	response := &v1beta1.AdmissionReview{}
	status := httpTest.Post(t, url, request, response, "")

	// Assert
	test.AssertEquals("", http.StatusOK, status, t)

	test.AssertFalse("Expected response to be non-nil", response.Response == nil, t)
	test.AssertEquals("", request.Request.UID, response.Response.UID, t)
	test.AssertFalse("", response.Response.Allowed, t)
	test.AssertFalse("Expected Result to be non-nil", response.Response.Result == nil, t)
	test.AssertEquals("", "Failure", response.Response.Result.Status, t)
	test.AssertEquals("", "Metadata field 'name' not found", response.Response.Result.Message, t)
	test.AssertEquals("", int32(http.StatusBadRequest), response.Response.Result.Code, t)
	test.AssertTrue("Expected PatchType to be nil", response.Response.PatchType == nil, t)
	test.AssertTrue("Expected Patch to be nil", response.Response.Patch == nil, t)

}

func TestMutateNamespace_when_error_parsing_raw_message(t *testing.T) {

	// Arrange
	mutator := createNamespaceMutatorForTesting()

	serverPort, appCtx := httpTest.StartTestServer([]base.Routes{mutator}, t)
	defer httpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/admissions/namespaces", serverPort)

	request := &v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			UID: types.UID("some-uid"),
			Object: runtime.RawExtension{
				Raw: json.RawMessage(`[]`),
			},
		},
	}

	// Act
	response := &v1beta1.AdmissionReview{}
	status := httpTest.Post(t, url, request, response, "")

	// Assert
	test.AssertEquals("", http.StatusOK, status, t)

	test.AssertFalse("Expected response to be non-nil", response.Response == nil, t)
	test.AssertEquals("", request.Request.UID, response.Response.UID, t)
	test.AssertFalse("", response.Response.Allowed, t)
	test.AssertFalse("Expected Result to be non-nil", response.Response.Result == nil, t)
	test.AssertEquals("", "Failure", response.Response.Result.Status, t)
	test.AssertEquals("", int32(http.StatusBadRequest), response.Response.Result.Code, t)
	test.AssertTrue("Expected PatchType to be nil", response.Response.PatchType == nil, t)
	test.AssertTrue("Expected Patch to be nil", response.Response.Patch == nil, t)

}

func TestMutateNamespace_when_malformed_message(t *testing.T) {

	// Arrange
	mutator := createNamespaceMutatorForTesting()

	serverPort, appCtx := httpTest.StartTestServer([]base.Routes{mutator}, t)
	defer httpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/admissions/namespaces", serverPort)

	request := &v1beta1.AdmissionReview{}

	// Act
	response := &httpUtils.ErrorMessage{}
	status := httpTest.Post(t, url, request, response, "")

	// Assert
	test.AssertEquals("", http.StatusBadRequest, status, t)
	test.AssertEquals("", http.StatusBadRequest, response.StatusCode, t)
	test.AssertEquals("", "Bad Request", response.Message, t)
	test.AssertEquals("", "[Request must be non-nil]", response.Detail, t)

}
