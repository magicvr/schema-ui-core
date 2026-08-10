package kernel

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// MigrationChecksum is the canonical ledger checksum for a migration: SHA-256
// (lower hex) of the normalized canonical SQL plus the data-transformer id
// (R6 C6.2 slice 3: module migration packages compute the same checksum the
// store ledger records).
func MigrationChecksum(stmts []string, transformID string) string {
	input := normalizeSQL(strings.Join(stmts, "\n")) + "\n" + transformID
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func normalizeSQL(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// This file implements the compiled-global Persistence catalog frozen in
// GOAL-005/attachments/r4-c1-freeze-package-draft.md §4. The only collection
// entry is Provider.CompiledPersistence(); the Registrar has no Persistence
// method, so migrations never enter the enablement-gated surface path.

// CollectPersistence calls CompiledPersistence() once per compiled provider and
// validates the merged catalog (freeze package §4.1):
//
//   - every migration carries a valid identity, version, name and checksum;
//   - global version, name and checksum are each unique; (module, name) is unique;
//   - versions are strictly consecutive (no gaps) within the collected set;
//   - tombstone and reconcile metadata are internally consistent;
//   - the result is deterministically ordered by version.
//
// Any violation fails closed; nothing is returned on error. Existing hardcoded
// store migrations (0001..0010) are current history, not the R4 terminal
// contract; module-owned migrations append global versions after them.
func CollectPersistence(providers []Provider) ([]MigrationContribution, error) {
	var catalog []MigrationContribution
	for _, provider := range providers {
		desc := provider.Descriptor()
		migrations, err := provider.CompiledPersistence()
		if err != nil {
			return nil, kernelError(CodeModuleInvalid, desc.ID, "collect persistence: %v", err)
		}
		for _, migration := range migrations {
			if migration.ModuleID == "" {
				migration.ModuleID = desc.ID
			}
			if err := validateMigration(migration.ModuleID, migration); err != nil {
				return nil, err
			}
			catalog = append(catalog, migration)
		}
	}
	return finalizePersistence(catalog)
}

func finalizePersistence(catalog []MigrationContribution) ([]MigrationContribution, error) {
	sort.Slice(catalog, func(i, j int) bool { return catalog[i].Version < catalog[j].Version })

	byVersion := map[int]string{}
	byName := map[string]string{}
	byChecksum := map[string]string{}
	byModuleName := map[string]string{}
	for _, m := range catalog {
		if previous, exists := byVersion[m.Version]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration version %d conflicts with %s", m.Version, previous)
		}
		byVersion[m.Version] = m.Name
		if previous, exists := byName[m.Name]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration name %q conflicts with %s", m.Name, previous)
		}
		byName[m.Name] = m.ModuleID
		if previous, exists := byChecksum[m.Checksum]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration checksum %q conflicts with %s", m.Checksum, previous)
		}
		byChecksum[m.Checksum] = m.ModuleID
		moduleName := m.ModuleID + "\x00" + m.Name
		if previous, exists := byModuleName[moduleName]; exists {
			return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration (%s, %s) conflicts with %s", m.ModuleID, m.Name, previous)
		}
		byModuleName[moduleName] = m.Name
		if m.ReconcileVersion != 0 {
			if m.ReconcileVersion <= 0 || m.ReconcileVersion > m.Version {
				return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration %q reconcile version %d must be in (0, %d]", m.Name, m.ReconcileVersion, m.Version)
			}
			if strings.TrimSpace(m.ReconcileChecksum) == "" {
				return nil, kernelError(CodeModuleInvalid, m.ModuleID, "migration %q reconcile version set without checksum", m.Name)
			}
		}
	}
	// No gaps: versions must be strictly consecutive.
	if len(catalog) > 1 {
		for i := 1; i < len(catalog); i++ {
			if catalog[i].Version != catalog[i-1].Version+1 {
				return nil, kernelError(CodeModuleInvalid, catalog[i].ModuleID, "migration version gap: %d followed by %d", catalog[i-1].Version, catalog[i].Version)
			}
		}
	}
	return catalog, nil
}
