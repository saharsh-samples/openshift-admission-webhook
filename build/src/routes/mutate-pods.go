package routes

import (
	"maw/integrations"
	"net/http"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

type podMutator struct {
	specialProvider integrations.SpecialProvider
	jsonUtils       utils.JSONUtils
}

// Register endpoint+method handlers
func (resource *podMutator) Register(agent base.RoutesAgent) {
	agent.RegisterPost("/admissions/pods", resource.MutatePod)
}

func (resource *podMutator) MutatePod(w http.ResponseWriter, r *http.Request) {

	h := &admissionHelper{jsonUtils: resource.jsonUtils, r: r, w: w}
	review, parseError := h.parseIncomingReview()
	if parseError != nil {
		return
	}

	// Determine if namespace is special
	isSpecial := resource.specialProvider.IsNamespaceSpecial(review.Request.Namespace)

	// if special, add node selector
	if isSpecial {
		h.allowWithPatches(
			review.Request.UID,
			[]string{`{ "op" : "add", "path": "/spec/tolerations/-", "value": { "key": "workload-type", "operator": "Equal", "value": "special", "effect": "NoSchedule" }}`},
		)
	} else {
		h.allowWithoutPatches(review.Request.UID)
	}

}
