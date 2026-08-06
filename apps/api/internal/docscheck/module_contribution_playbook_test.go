// Package docscheck verifies documentation contracts that product playbooks
// claim against the real repository layout (VP-004 / workspace-004 deliverable).
package docscheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// test file lives at apps/api/internal/docscheck → repo root is ../../../../
	root := filepath.Clean(filepath.Join(wd, "..", "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "docs", "architecture", "module-contribution-playbook.md")); err != nil {
		// Fallback: walk up looking for go.mod of apps/api then one more up.
		dir := wd
		for i := 0; i < 8; i++ {
			candidate := filepath.Join(dir, "docs", "architecture", "module-contribution-playbook.md")
			if _, err := os.Stat(candidate); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
		t.Fatalf("cannot locate repo root from %s: playbook missing (%v)", wd, err)
	}
	return root
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func TestModuleContributionPlaybookShippedContent(t *testing.T) {
	root := repoRoot(t)
	playbookPath := filepath.Join(root, "docs", "architecture", "module-contribution-playbook.md")
	body := read(t, playbookPath)

	requiredSnippets := []string{
		"## 1. 新增一方标准 Admin 功能模块 · 必须完成（MUST）",
		"## 2. 明确不必做 / 禁止做（DO NOT）",
		"## 3. Core vs 模块责任 · 归属判定",
		"模块 id / 版本 / 内核 API 范围 / 依赖",
		"核心六项贡献",
		"组合根静态候选注册",
		"Profile / `modules.enabled`",
		"全局迁移台账",
		"验证 / 回归最小集",
		"不要为接模块改 Renderer",
		"不要在生产路径静默使用静态 Manifest",
		"不要私建平行认证",
		"不要把「按需能力」当成核心六项",
		"不要做运行时插件",
		"apps/api/internal/modules/users",
		"apps/api/internal/composition/composition.go",
		"apps/api/internal/kernel/profile.go",
		"apps/api/internal/modules/compiled/persistence.go",
	}
	for _, snip := range requiredSnippets {
		if !strings.Contains(body, snip) {
			t.Errorf("playbook missing required snippet: %q", snip)
		}
	}

	// Cited real paths must exist on disk (not invented).
	for _, rel := range []string{
		"apps/api/internal/modules/users/provider.go",
		"apps/api/internal/composition/composition.go",
		"apps/api/internal/kernel/profile.go",
		"apps/api/internal/modules/compiled/persistence.go",
		"docs/architecture/module-architecture.md",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Errorf("cited path missing: %s (%v)", rel, err)
		}
	}
}

func TestPlaybookDiscoverabilityFromOverviewAndQuickstart(t *testing.T) {
	root := repoRoot(t)
	linkToken := "module-contribution-playbook.md"

	overview := read(t, filepath.Join(root, "docs", "architecture", "overview.md"))
	if !strings.Contains(overview, linkToken) {
		t.Fatal("overview.md does not link to module-contribution-playbook.md")
	}

	quickstart := read(t, filepath.Join(root, "QUICKSTART.md"))
	if !strings.Contains(quickstart, linkToken) {
		t.Fatal("QUICKSTART.md does not link to module-contribution-playbook.md")
	}

	arch := read(t, filepath.Join(root, "docs", "architecture", "module-architecture.md"))
	if !strings.Contains(arch, linkToken) {
		t.Fatal("module-architecture.md does not link out to playbook")
	}

	// AGENTS/Skills must not be required as the default discovery path:
	// playbook itself states default sufficiency via overview+QUICKSTART.
	playbook := read(t, filepath.Join(root, "docs", "architecture", "module-contribution-playbook.md"))
	if !strings.Contains(playbook, "不要求改 `AGENTS.md`") && !strings.Contains(playbook, "不要求改 AGENTS") {
		// Accept either backtick form used in the playbook §4
		if !strings.Contains(playbook, "AGENTS.md") || !strings.Contains(playbook, "Skills") {
			t.Fatal("playbook should document AGENTS/Skills non-requirement for AI discovery")
		}
	}
}

func TestWorkspace004RootCloseoutArtifacts(t *testing.T) {
	root := repoRoot(t)
	base := filepath.Join(root, "docs", "workspace-004-module-contribution-readiness", "GOAL-001-module-contribution-readiness")
	meta := read(t, filepath.Join(base, "00-meta.md"))
	if !strings.Contains(meta, "status: done") {
		t.Fatal("Root 00-meta status is not done")
	}
	if !strings.Contains(meta, "progress: 4/4") {
		t.Fatal("Root 00-meta progress is not 4/4")
	}
	for _, mark := range []string{"[x] **S1**", "[x] **S2**", "[x] **S3**", "[x] **S4**"} {
		if !strings.Contains(meta, mark) {
			t.Errorf("Root missing checked checkpoint: %s", mark)
		}
	}
	if !strings.Contains(meta, "I-001") || !strings.Contains(meta, "verified") {
		t.Fatal("I-001 should be verified in Root meta")
	}

	tree := read(t, filepath.Join(root, "docs", "workspace-004-module-contribution-readiness", "goal-tree.md"))
	if !strings.Contains(tree, "done") || !strings.Contains(tree, "4/4") {
		t.Fatal("goal-tree does not show Root done 4/4")
	}

	audit := read(t, filepath.Join(base, "03-audit", "A-001-root-closeout-self.md"))
	if !strings.Contains(audit, "verdict: pass") && !strings.Contains(audit, "verdict | **pass**") && !strings.Contains(audit, "**pass**") {
		t.Fatal("A-001 missing pass verdict")
	}
	if !strings.Contains(strings.ToLower(audit), "source") || !strings.Contains(audit, "self") {
		t.Fatal("A-001 should be source self")
	}

	// principles.md must not have been rewritten as the product playbook target
	// (content authority stays under architecture product docs).
	if _, err := os.Stat(filepath.Join(root, "docs", "architecture", "principles.md")); err != nil {
		t.Fatal("principles.md missing unexpectedly")
	}
}
