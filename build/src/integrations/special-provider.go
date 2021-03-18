package integrations

// SpecialProvider for metadata related to special compliance
type SpecialProvider interface {
	IsNamespaceSpecial(namespace string) bool
}

type specialProvider struct{}

func (provider *specialProvider) IsNamespaceSpecial(namespace string) bool {
	return namespace == "special-namespace"
}
