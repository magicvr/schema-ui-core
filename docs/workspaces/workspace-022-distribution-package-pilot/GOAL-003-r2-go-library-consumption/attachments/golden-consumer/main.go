// 黄金下游仓 · 最小真实组合根（R2 装配闭环第一验证 · S3 前奏）
//
// 仅 `go get`（本地 replace）+ 自建组合根的消费形态：
//   - import A 层（apps/api/kernel）与 B 层模块（apps/api/modules/dashboard）
//   - ResolveProfile → NewRegistry(BuiltinModules) → dashboard 提供者装配
package main

import (
	"fmt"
	"os"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/dashboard"
)

func main() {
	plan, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve:", err)
		os.Exit(1)
	}
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		fmt.Fprintln(os.Stderr, "registry:", err)
		os.Exit(1)
	}
	_ = registry

	d := dashboard.New()
	desc := d.Descriptor()
	fmt.Printf("kernel=%s profile=%s dashboard=%s@%s range=%s modules=%d\n",
		kernel.KernelAPIVersion, plan.Name, desc.ID, desc.Version, desc.KernelAPIRange, len(plan.Modules))
}