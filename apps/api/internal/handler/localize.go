package handler

// Error localization writer (VP-007 S4 · I-L10N-004 path a).
// See internal/errorcatalog for the catalog and negotiation semantics.

import (
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
)

// writeLocalizedError writes the compatible error envelope with server-side
// locale negotiation. Cataloged codes get a localized message + messageKey +
// Content-Language; uncataloged codes keep the English generic message with no
// messageKey (never leaks diagnostics).
func writeLocalizedError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	locale := errorcatalog.Negotiate(r)
	body, contentLanguage, cataloged := errorcatalog.Body(code, message, locale)
	if !cataloged {
		writeJSON(w, status, map[string]string{"error": code, "message": message})
		return
	}
	w.Header().Set("Content-Language", contentLanguage)
	writeJSON(w, status, body)
}
