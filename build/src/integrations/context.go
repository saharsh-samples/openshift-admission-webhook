package integrations

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	// no dependencies
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	SpecialProvider SpecialProvider
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {
	out := &ContextOut{}
	out.SpecialProvider = &specialProvider{}
	return out
}
