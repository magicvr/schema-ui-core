// 实验桩（E-001）：验证 A 层（internal/kernel）是否可被外部模块 import。
// 预期：Go internal 规则拒绝（"use of internal package ... not allowed"）。
package main

import (
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

func main() {
	fmt.Println(kernel.KernelAPIVersion)
}