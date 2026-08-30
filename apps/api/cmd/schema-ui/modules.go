package main

import (
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// builtinModules 列出可用标准模块（CLI 与 kernel 同模块，直接读取；
// 保证清单与运行时装配语义一致）。
func builtinModules() []moduleInfo {
	mods := kernel.BuiltinModules()
	out := make([]moduleInfo, 0, len(mods))
	for _, m := range mods {
		out = append(out, moduleInfo{id: m.ID, version: m.Version, kernelAPI: m.KernelAPIRange})
	}
	return out
}