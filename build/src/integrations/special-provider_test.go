package integrations

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestIsNamespaceSpecial(t *testing.T) {

	provider := Bootstrap(nil).SpecialProvider

	test.AssertTrue("Expected namespace to be special", provider.IsNamespaceSpecial("special-namespace"), t)
	test.AssertFalse("Expected namespace to be special", provider.IsNamespaceSpecial("ordinary-namespace"), t)
	test.AssertFalse("Expected namespace to be special", provider.IsNamespaceSpecial("ignored-namespace"), t)
}
