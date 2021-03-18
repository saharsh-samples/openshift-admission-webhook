package routes

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"

	"github.com/saharsh-samples/go-mux-sql-starter/http/utils"
)

func TestBootstrap(t *testing.T) {

	// Arrange
	in := &ContextIn{
		JSONUtils: utils.Bootstrap(nil).JSONUtils,
	}

	// Act
	out := Bootstrap(in)

	// Assert
	test.AssertEquals("", 3, len(out.RoutesToRegister), t)
}
