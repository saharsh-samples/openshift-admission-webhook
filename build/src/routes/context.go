package routes

import (
	"maw/integrations"

	"github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/http/routes"
	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

// ContextIn describes dependecies needed by this package
type ContextIn struct {

	// HTTP Utils dependencies
	JSONUtils utils.JSONUtils

	// integrations
	SpecialProvider integrations.SpecialProvider
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	RoutesToRegister []http.Routes
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	out := &ContextOut{}

	// routes
	out.RoutesToRegister = []http.Routes{

		// health check
		&routes.LivenessCheck{},

		// mutating admission webhook for pods
		&podMutator{
			specialProvider: in.SpecialProvider,
			jsonUtils:       in.JSONUtils,
		},

		// mutating admission webhook for namespaces
		&namespaceMutator{
			specialProvider: in.SpecialProvider,
			jsonUtils:       in.JSONUtils,
		},
	}

	return out
}
