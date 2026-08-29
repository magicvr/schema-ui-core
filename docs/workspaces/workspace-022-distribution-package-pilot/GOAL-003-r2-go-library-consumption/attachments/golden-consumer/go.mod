module golden.local/consumer

go 1.26

require github.com/magicvr/schema-ui-core/apps/api v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/image v0.45.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// 试点黄金消费仓：R2 装配闭环验证用（嵌套 go.mod，主仓 ./... 自动排除）。
// replace 指向本仓 apps/api；发布流水线（R5）落成前不依赖真实 tag。
replace github.com/magicvr/schema-ui-core/apps/api => ../../../../../../apps/api
