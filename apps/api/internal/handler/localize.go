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

// writeLocalizedFieldError is writeLocalizedError plus a fieldErrors array
// (GOAL-014 D-002 §2.1): field-level validation failures (create/patch) carry
// the concrete {field, reason} pairs so the host can inline them on inputs.
// The envelope stays backward compatible (fieldErrors is an optional key).
func writeLocalizedFieldError(w http.ResponseWriter, r *http.Request, status int, code, message string, fields []errorcatalog.FieldError) {
	locale := errorcatalog.Negotiate(r)
	body, contentLanguage, cataloged := errorcatalog.BodyWithFields(code, message, locale, fields)
	if !cataloged {
		writeJSON(w, status, map[string]any{"error": code, "message": message, "fieldErrors": fields})
		return
	}
	w.Header().Set("Content-Language", contentLanguage)
	writeJSON(w, status, body)
}
