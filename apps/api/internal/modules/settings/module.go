package settings

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// Register mounts the Settings module's HTTP contribution. The module owns
// the registration boundary; the handler package only supplies the stable
// protocol implementation during this bounded migration slice.
func Register(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	handler.RegisterSettings(mux, a, st)
}
