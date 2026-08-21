package errorcatalog

import (
	"encoding/json"
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

// WriteLocalizedError writes the shared error envelope used by both the auth
// middleware and ordinary handlers. Keeping the wire writer here prevents the
// two HTTP boundaries from drifting in localization or correlation behavior.
func WriteLocalizedError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	locale := Negotiate(r)
	body, contentLanguage, cataloged := Body(code, message, locale)
	if id := requestid.FromContext(r.Context()); id != "" {
		if body == nil {
			body = map[string]any{}
		}
		body[requestid.BodyName] = id
	}
	if !cataloged {
		body = map[string]any{"error": code, "message": message}
		if id := requestid.FromContext(r.Context()); id != "" {
			body[requestid.BodyName] = id
		}
	} else {
		w.Header().Set("Content-Language", contentLanguage)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
