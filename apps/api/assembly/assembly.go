// Package assembly 提供公开装配工厂（B+ 层 · 契约冻结面 §2b 候选）。
//
// 试点（VP-022 R2 S3 · D-003 方案 β · 2026-08-29 用户裁决）最小面：
// 下游组合根经本包注入基架实现（store / auth / mail），无需命名 internal 类型
// ——Go 类型推断对 *auth.Authenticator、kernel.Store 等返回值直接消费。
//
// 版本语义：0.1.0 experimental；signature 冻结于契约冻结面 v1.1.0（增列 B+ 层）。
// 本包同模块树，可合法 import apps/api/internal/*（C 层）。
package assembly

import (
	"context"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
)

// OpenStore 打开内核 Store（SQLite / PG 双方言；迁移随 Open 自动 apply，
// 全局 checksum 台账校验在内部完成）。
//
//	catalog = 模块迁移贡献（如 compiled.PersistenceCatalog()）；nil = 仅内核台账。
func OpenStore(ctx context.Context, dialect kernel.Dialect, path, dsn string, catalog []kernel.MigrationContribution) (kernel.Store, error) {
	return store.Open(ctx, store.OpenOptions{
		Dialect: dialect,
		Path:    path,
		DSN:     dsn,
	}, catalog)
}

// NewAuthenticator 构造 JWT 认证器（current 单密钥形态）。
//
//	runner 通常直接传 OpenStore 的返回值——kernel.Store 自动满足 authsession.TxRunner。
func NewAuthenticator(secret []byte, accessTTL, refreshTTL time.Duration, runner authsession.TxRunner) *auth.Authenticator {
	return auth.New(secret, accessTTL, refreshTTL, runner, false)
}

// NewMailSender 构造站内出站记录渠道（mock 默认形态；管理员面可检视报文）。
func NewMailSender(st kernel.Store, retentionCap int) kernel.MailSender {
	return mail.NewOutboxSink(st, retentionCap)
}