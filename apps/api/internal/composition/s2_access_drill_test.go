package composition

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// S2 非领域化接入演练（VP-008 §最小可枚举证据面 6）：
// 验证一个新增 standard-admin 风格 Probe 模块通过 kernel.Provider 贡献六面
// （HTTP / Schema / Authorization / Navigation / Manifest），在 composition root
// 装配后，其路由、Schema 文档、权限、导航与 Manifest 页面全部经聚合表面发布，
// 且不修改 Renderer/Shell 中央业务注册（central registration 无 probe 分支）。
//
// Probe 只存在于 test-only 候选集，绝不进入 mvp/admin 默认 Profile 或生产 Manifest。

// probeProvider 是 test-only 的「非领域化」probe：声明一个 generic 资源
// (generic.items) 的读+写表面，映射为后续业务模块的框架共性能力
// （list/detail/create/update/delete + read/write 权限 + 导航 + Manifest fragment）。
type probeProvider struct {
	desc kernel.Module
}

func (p *probeProvider) Descriptor() kernel.Module { return p.desc }

func (p *probeProvider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	// Probe 无自己的 schema 迁移：归属 core.auth-session（与 admin.users 同模式 C3.3）。
	return nil, nil
}

func (p *probeProvider) Register(ctx context.Context, reg kernel.Registrar) error {
	id := p.desc.ID
	// M2 六项：HTTP
	if err := reg.HTTP(kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: id, Key: "GET /api/probe-items"},
		Method:               "GET",
		Pattern:              "/api/probe-items",
		Handler:              http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
	}); err != nil {
		return err
	}
	// Schema
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: id, Key: "probe-items"},
		PageID:               "probe-items",
		Resources:            []string{"probe-items"},
		Actions:              []string{"list", "create", "detail", "update", "delete"},
		DataSource:           "/api/probe-items",
		Owner:                id,
		Document:             probeDocument("probe-items"),
	}); err != nil {
		return err
	}
	// Authorization
	if err := reg.Authorization(kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: id, Key: "probe.read"},
		Permission:           "probe.read",
		Resource:             "probe-items",
		Action:               "read",
		PolicyID:             "system.admin-editor-viewer",
		SystemDataVersion:    1,
	}); err != nil {
		return err
	}
	// Navigation
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: id, Key: "menu_probe"},
		NodeID:               "menu_probe",
		PageID:               "probe-items",
		Order:                99,
		Label:                "Probe Items",
		Visibility:           "system.admin",
		Permission:           "probe.read",
		SystemDataVersion:    1,
	}); err != nil {
		return err
	}
	// Manifest fragment（ProtocolVersion/RequiredCapabilities/app 与现行 standard-admin 一致）
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: id, Key: "probe"},
		FragmentID:           "probe",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON: []byte(`{
  "protocolVersion": "2.7",
  "requiredCapabilities": ["app.manifest", "app.navigation"],
  "app": {
    "appId": "schema-ui-core",
    "name": "Schema UI Core",
    "description": "The schema-ui-core administration workspace."
  },
  "pages": [
    {
      "pageId": "probe-items",
      "title": "Probe Items",
      "schemaUrl": "/api/schema/probe-items",
      "route": "/probe-items",
      "titleKey": "manifest.title.probe"
    }
  ],
  "navigation": {
    "sidebar": [
      {
        "pageRef": "probe-items",
        "label": "Probe Items",
        "icon": "box",
        "visibleWhen": {"when": "$context.features.menu_probe == true"},
        "labelKey": "manifest.nav.probe"
      }
    ]
  }
}`),
	})
}

func probeDocument(pageID string) []byte {
	raw, _ := json.Marshal(map[string]any{
		"meta": map[string]any{"pageId": pageID},
		"body": map[string]any{},
	})
	return raw
}

func probeModule() kernel.Module {
	return kernel.Module{
		ID:             "admin.probe",
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes:      []string{"GET /api/probe-items"},
			Pages:       []string{"probe-items"},
			Navigation:  []string{"menu_probe"},
			Permissions: []string{"probe.read"},
			Fragments:   []string{"probe"},
		},
	}
}

// TestS2AccessDrill_ProbeModuleSurfacesThroughComposition 验证：
// (1) 把 test-only probe 模块并入编译候选（BuiltinModules + probe），
// (2) 在 admin plan 中启用 probe，
// (3) composition root 装配后，probe 的 route / schema / manifest 均经聚合表面发布，
//
//	且 central registration（Web Renderer/Shell）无 probe 分支 —— 见下方断言
//	「probe 不依赖任何中央业务注册改动」。
func TestS2AccessDrill_ProbeModuleSurfacesThroughComposition(t *testing.T) {
	all := append([]kernel.Module(nil), kernel.BuiltinModules()...)
	all = append(all, probeModule())
	registry, err := kernel.NewRegistry(all)
	if err != nil {
		t.Fatalf("registry with probe: %v", err)
	}

	// 用 admin profile 作为宿主，显式并入 probe（探针只进 test-only 候选集）。
	adminResolution, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatalf("resolve admin: %v", err)
	}
	plan, err := registry.Resolve(append(adminResolution.Modules, probeModule().ID))
	if err != nil {
		t.Fatalf("plan w/ probe: %v", err)
	}

	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	// W13 F-010: schemas are authenticated now — inject the dev session so the
	// probe schema assertion reaches the handler.
	a := auth.New([]byte("test-secret"), 0, 0, st, true)

	// composition root 装配：extra 传入 probe（等效 M3 静态 import + plan 分支）。
	// 不包含任何 handler / Web Renderer-Shell 中央业务注册改动。
	probe := &probeProvider{desc: probeModule()}
	jobRuntime, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	mux, err := newMuxWithExtraProviders(
		&config.Config{ProfileName: "admin"},
		a,
		st,
		authsession.NewRepository(st),
		operationlog.NewRepository(st),
		settingsrepository.New(st),
		plan,
		&readinessGate{},
		jwtSecret("test-secret"),
		jobRuntime,
		[]kernel.Provider{probe},
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	if err != nil {
		t.Fatalf("newMuxWithExtraProviders with probe: %v", err)
	}

	// (a) probe route 可访问
	assertStatus(t, mux, "/api/probe-items", http.StatusOK)

	// (b) probe schema 文档可经 /api/schema/{pageId} 获取
	assertStatus(t, mux, "/api/schema/probe-items", http.StatusOK)

	// (c) probe 页面进入聚合 Manifest
	manifest := httptest.NewRecorder()
	mux.ServeHTTP(manifest, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if manifest.Code != http.StatusOK {
		t.Fatalf("manifest status = %d", manifest.Code)
	}
	var doc struct {
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(manifest.Body.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	pageIDs := map[string]bool{}
	for _, p := range doc.Pages {
		pageIDs[p.PageID] = true
	}
	if !pageIDs["probe-items"] {
		t.Fatalf("probe page missing from manifest: %v", pageIDs)
	}
}

// TestS2AccessDrill_ProbeAbsentFromDefaultProfiles 验证：
// probe / 业务探针绝不进入 mvp/admin 默认启用集（Part 6：probe default 只进
// test-only 候选集）。
func TestS2AccessDrill_ProbeAbsentFromDefaultProfiles(t *testing.T) {
	for _, name := range []kernel.ProfileName{kernel.ProfileMVP, kernel.ProfileAdmin} {
		resolution, err := kernel.ResolveProfile(string(name), nil)
		if err != nil {
			t.Fatalf("%s resolve: %v", name, err)
		}
		for _, id := range resolution.Modules {
			if id == "admin.probe" || id == "probe" {
				t.Fatalf("%s default profile must not include probe module", name)
			}
		}
	}
}

func assertStatus(t *testing.T, mux *http.ServeMux, path string, want int) {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if rec.Code != want {
		t.Fatalf("GET %s = %d, want %d; body=%s", path, rec.Code, want, rec.Body.String())
	}
}
