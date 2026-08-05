package systemdata

// Stable policy identifiers are part of the contribution contract. They map a
// finalized permission/navigation contribution to the system roles that own
// its default grant; custom roles are never touched by reconciliation.
const (
	PolicyAdmin             = "system.admin"
	PolicyAdminEditor       = "system.admin-editor"
	PolicyAdminEditorViewer = "system.admin-editor-viewer"
	SystemDataVersion       = 1
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
