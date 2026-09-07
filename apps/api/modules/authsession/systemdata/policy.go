package systemdata

// Stable policy identifiers are part of the contribution contract. They map a
// finalized permission/navigation contribution to the system roles that own
// its default grant; custom roles are never touched by reconciliation.
const (
	PolicyAdmin             = "system.admin"
	PolicyAdminEditor       = "system.admin-editor"
	PolicyAdminEditorViewer = "system.admin-editor-viewer"
	// SystemDataVersion stamps every permission/navigation contribution
	// checksum. Bump it whenever the CONTENT of any contribution changes
	// (label/order/policy/…); the reconcile ledger accepts the new checksums
	// only on a version increase (same-version mismatch = tamper → fail closed,
	// see checkLedger). v1 → v2 (2026-09-06 · workspace-031 post-closure):
	// digital-offer nav label "Digital offers" → "Digital products" and the
	// new menu_digitaloffer_purchases entry (E-013/E-014).
	SystemDataVersion = 2
)

func rolesForPolicy(policyID string) ([]string, bool) {
	switch policyID {
	case PolicyAdmin:
		return []string{"admin"}, true
	case PolicyAdminEditor:
		return []string{"admin", "editor"}, true
	case PolicyAdminEditorViewer:
		return []string{"admin", "editor", "viewer"}, true
	default:
		return nil, false
	}
}
