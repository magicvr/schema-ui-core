package obs

import "net/http"

// InstrumentedMux wraps a ServeMux so every registration is measured under
// its owning module_id (R2 D-001 §1/§2). It is the single interception point
// for the composition mux: central handler registrations (mux.Handle and
// mux.HandleFunc) and module-contributed routes all flow through it.
//
// Ownership: patterns declared via Own before Handle use that module id;
// everything else defaults to ModuleIDCore. The embedded ServeMux is what is
// ultimately mounted, so matching semantics are unchanged.
type InstrumentedMux struct {
	*http.ServeMux
	observer *Observer
	owners   map[string]string
}

// NewInstrumentedMux builds the wrapper. observer may be nil (metrics
// disabled): registrations then pass through untouched.
func NewInstrumentedMux(observer *Observer) *InstrumentedMux {
	return &InstrumentedMux{
		ServeMux: http.NewServeMux(),
		observer: observer,
		owners:   map[string]string{},
	}
}

// Own declares the contributing module for a full registration pattern
// ("GET /api/users"). Empty module ids fall back to ModuleIDCore.
func (m *InstrumentedMux) Own(pattern, moduleID string) {
	if moduleID == "" {
		moduleID = ModuleIDCore
	}
	m.owners[pattern] = moduleID
}

func (m *InstrumentedMux) ownerOf(pattern string) string {
	if id, ok := m.owners[pattern]; ok {
		return id
	}
	return ModuleIDCore
}

// Handle registers the handler, wrapped with metrics instrumentation when an
// Observer is present.
func (m *InstrumentedMux) Handle(pattern string, handler http.Handler) {
	m.ServeMux.Handle(pattern, m.observer.Wrap(pattern, m.ownerOf(pattern), handler))
}

// HandleFunc mirrors http.ServeMux.HandleFunc so callers using the function
// form stay instrumented.
func (m *InstrumentedMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(handler))
}
