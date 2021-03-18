package context

import (
	"fmt"
	"maw/config"
	"maw/integrations"
	"maw/routes"

	"github.com/saharsh-samples/go-mux-sql-starter/app"
	"github.com/saharsh-samples/go-mux-sql-starter/http"
	httpUtils "github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

// Build bootstraps all app submodules and creates overall app context
func Build() *app.ContextOut {

	// app config
	appConfig := config.AppConfig{}
	appConfig.InitFromEnvironment()

	// http utils
	httpUtilsOut := httpUtils.Bootstrap(&httpUtils.ContextIn{})

	// integrations
	integrationsOut := integrations.Bootstrap(nil)

	// http routes
	routesOut := routes.Bootstrap(&routes.ContextIn{

		// http utils
		JSONUtils: httpUtilsOut.JSONUtils,

		// integrations
		SpecialProvider: integrationsOut.SpecialProvider,
	})

	// http base
	httpOut := http.Bootstrap(&http.ContextIn{
		Port:             appConfig.HTTPConfig.Port,
		RoutesToRegister: routesOut.RoutesToRegister,
		TLSConfiguration: &http.TLSConfiguration{
			CertFile: appConfig.HTTPConfig.TLSCertFilePath,
			KeyFile:  appConfig.HTTPConfig.TLSKeyFilePath,
		},
	})

	// app
	return app.Bootstrap(&app.ContextIn{

		// server
		StartupTimeoutInSeconds: appConfig.HTTPConfig.StartupTimeoutInSeconds,
		HTTPServer:              httpOut.Server,

		// Add all shutdown hooks here
		ShutdownHooks: []app.ShutdownHook{
			func() {
				// make sure this is at beginning (runs last)
				fmt.Println("Graceful shutdown of application complete.")
			},
		},
	})
}
