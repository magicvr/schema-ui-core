module golden.local/consumer

go 1.26

require github.com/magicvr/schema-ui-core/apps/api v0.0.0

// 试点黄金消费仓：R2 装配闭环验证用（嵌套 go.mod，主仓 ./... 自动排除）。
// replace 指向本仓 apps/api；发布流水线（R5）落成前不依赖真实 tag。
replace github.com/magicvr/schema-ui-core/apps/api => ../../../../../../apps/api