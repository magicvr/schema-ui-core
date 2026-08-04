package activity

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// Register mounts the Activity module's HTTP contribution. Operation-log
// persistence remains a core capability even when this read-only surface is
// not selected.
func Register(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	handler.RegisterActivity(mux, a, st)
}
