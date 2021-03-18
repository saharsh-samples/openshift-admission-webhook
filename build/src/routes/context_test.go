package routes

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"

	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

func TestBootstrap(t *testing.T) {

	// Arrange
	utilsOut := utils.Bootstrap(nil)
	in := &ContextIn{
		JSONUtils: utilsOut.JSONUtils,
		URLUtils:  utilsOut.URLUtils,
	}

	// Act
	out := Bootstrap(in)

	// Assert
	test.AssertEquals("", 1, len(out.RoutesToRegister), t)
}
