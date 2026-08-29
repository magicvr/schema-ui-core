// 黄金下游仓 · 完整装配闭环（R2 S3 · 方案 β 实证）
//
// 外部模块（golden.local/consumer）消费形态：
//   - A 层：apps/api/kernel
//   - B 层：apps/api/modules/{authsession,operationlog,users,compiled}
//   - B+ 层：apps/api/assembly（公开装配工厂——类型推断消费，无 internal 命名）
//
// 流程：迁移台账收集 → OpenStore（SQLite，迁移随 Open apply）→ 认证/仓库/邮件装配
// → users Provider → Profile 解析 → Registry 校验 → RegisterContributions。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/assembly"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/users"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: consumer <sqlite-path>")
	}
	dbPath := os.Args[1]
	ctx := context.Background()

	// 1. 迁移台账收集（B 层 compiled：全部候选迁移，checksum 校验在 store.Open 内）
	catalog, err := compiled.PersistenceCatalog()
	if err != nil {
		log.Fatal("catalog:", err)
	}

	// 2. 打开双方言之一：SQLite（内核 Store 端口）
	st, err := assembly.OpenStore(ctx, kernel.Dialect("sqlite"), dbPath, "", catalog)
	if err != nil {
		log.Fatal("store:", err)
	}
	defer st.Close()

	// 3. 基架装配（类型推断消费：*auth.Authenticator 不命名）
	authn := assembly.NewAuthenticator([]byte("consumer-test-secret"), 15*time.Minute, 24*time.Hour, st)
	repo := authsession.NewRepository(st)
	ops := operationlog.NewRepository(st)
	mailer := assembly.NewMailSender(st, 1000)

	// 4. 标准模块 users 全链装配（B 层公开构造）
	usersProvider := users.New(authn, repo, ops, mailer, "http://localhost")

	// 5. Profile 解析 + Registry 校验（fail-closed）
	plan, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		log.Fatal("resolve:", err)
	}
	builtin := kernel.BuiltinModules()
	if _, err := kernel.NewRegistry(builtin); err != nil {
		log.Fatal("registry:", err)
	}

	// 6. 贡献注册（users 在 admin profile 内 → 全部贡献进入集合）
	byID := make(map[string]kernel.Module, len(builtin))
	for _, m := range builtin {
		byID[m.ID] = m
	}
	planModules := make([]kernel.Module, 0, len(plan.Modules))
	for _, id := range plan.Modules {
		if m, ok := byID[id]; ok {
			planModules = append(planModules, m)
		}
	}
	contribs, err := kernel.RegisterContributions(ctx, kernel.Plan{
		Modules:      planModules,
		Capabilities: kernel.StandardAdminCapabilities(),
	}, []kernel.Provider{usersProvider})
	if err != nil {
		log.Fatal("contribs:", err)
	}

	fmt.Printf("kernel=%s profile=%s dialect=%s fresh=%v "+
		"contrib{routes=%d pages=%d perms=%d nav=%d frag=%d} "+
		"users_module=%s\n",
		kernel.KernelAPIVersion, plan.Name, st.Dialect(), st.WasFresh(),
		len(contribs.Routes), len(contribs.Pages), len(contribs.Permissions),
		len(contribs.Navigation), len(contribs.Fragments),
		users.ModuleID)
}