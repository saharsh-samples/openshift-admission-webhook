package routes

import (
	"maw/integrations"
	"net/http"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

type namespaceMutator struct {
	specialProvider integrations.SpecialProvider
	jsonUtils       utils.JSONUtils
}

// Register endpoint+method handlers
func (resource *namespaceMutator) Register(agent base.RoutesAgent) {
	agent.RegisterPost("/admissions/namespaces", resource.MutateNamespace)
}

func (resource *namespaceMutator) MutateNamespace(w http.ResponseWriter, r *http.Request) {

	h := &admissionHelper{jsonUtils: resource.jsonUtils, r: r, w: w}
	review, parseError := h.parseIncomingReview()
	if parseError != nil {
		return
	}

	// Determine if namespace is special
	nsName, err := getResourceMetadataFieldAsString(review, "name")
	if err != nil {
		h.denyAdmission(review.Request.UID, err.Error(), http.StatusBadRequest)
		return
	}
	isSpecial := resource.specialProvider.IsNamespaceSpecial(nsName)

	// if special, add node selector
	if isSpecial {
		h.allowWithPatches(
			review.Request.UID,
			[]string{`{ "op" : "add", "path": "/metadata/annotations/openshift.io~1node-selector", "value": "workload-type=special"}`},
			resource.jsonUtils,
			w,
		)
	} else {
		h.allowWithoutPatches(review.Request.UID, resource.jsonUtils, w)
	}

}
