package routes

import (
	"encoding/json"
	"fmt"
	"maw/model"
	"net/http"

	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getRawRequest(review *model.AdmissionReview) (map[string]interface{}, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(review.Request.Object.Raw, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getResourceMetadataFieldAsString(review *model.AdmissionReview, metadataField string) (string, error) {

	reqObj, err := getRawRequest(review)
	if err != nil {
		return "", err
	}

	if metadata, found := (reqObj["metadata"]).(map[string]interface{}); found {
		if val, found := (metadata[metadataField]).(string); found {
			return val, nil
		}
	}

	return "", fmt.Errorf("Metadata field '%v' not found", metadataField)
}

// ---
// STATEFUL ADMISSION HELPER - should be created on a per-request basis
// ---

type admissionHelper struct {
	jsonUtils utils.JSONUtils
	r         *http.Request
	w         http.ResponseWriter
}

func (h *admissionHelper) parseIncomingReview() (*model.AdmissionReview, error) {

	// parse request
	review := &model.AdmissionReview{}
	parseError := h.jsonUtils.ParseJSONRequest(h.r, review, h.w)
	if parseError != nil {
		// ParseJSONRequest does proper HTTP-level error handling
		return nil, parseError
	}

	fmt.Println("Received admission review request with UID", review.Request.UID)
	return review, nil

}

func (h *admissionHelper) denyAdmission(requestUID types.UID, reason string, statusCode int) {
	h.jsonUtils.SetJSONResponse(h.w, http.StatusOK, &v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     requestUID,
			Allowed: false,
			Result: &v1.Status{
				Status:  "Failure",
				Message: reason,
				Code:    int32(statusCode),
			},
		},
	})
}

func (h *admissionHelper) allowWithoutPatches(requestUID types.UID) {
	h.jsonUtils.SetJSONResponse(h.w, http.StatusOK, &v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     requestUID,
			Allowed: true,
		},
	})
}

func (h *admissionHelper) allowWithPatches(requestUID types.UID, patches []string) {

	// create response object
	response := &v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     requestUID,
			Allowed: true,
		},
	}

	// add patches to response
	numOfPatches := len(patches)
	if numOfPatches > 0 {

		patchString := "["

		for i, patch := range patches {
			patchString = fmt.Sprintf("%v %v", patchString, patch)
			if i != numOfPatches-1 {
				patchString = fmt.Sprintf("%v,", patchString)
			}
		}

		patchString = fmt.Sprintf("%v ]", patchString)

		response.Response.PatchType = func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}()
		response.Response.Patch = json.RawMessage(patchString)
	}

	// respond
	h.jsonUtils.SetJSONResponse(h.w, http.StatusOK, response)
}
