package handler

import "net/http"

// routeRegistrar is the mux surface registration helpers mount routes on.
// The composition root passes an instrumentation wrapper (VP-015 R2 /
// workspace-015 GOAL-003) implementing the same two methods as
// *http.ServeMux, so central registrations are measured like module
// contributions; plain *http.ServeMux keeps working unchanged.
type routeRegistrar interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
