// 实验桩（E-003）：量产模块 users 的外部装配可行性。
// 目标：证实现有 B 层模块（users）的依赖链是否可被外部模块完整消费。
// 方法：编译期引用 users 包，观察依赖链加载结果（预期 internal 规则阻断）。
package main

import (
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/users"
)

var _ = kernel.KernelAPIVersion
var _ = users.ModuleID