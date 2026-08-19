package handler

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

func auditDetail(action string, after map[string]any) *string {
	raw, err := operationlog.NewDetail(action, nil, after)
	if err != nil {
		slog.Error("operation log detail", "action", action, "err", err)
		return nil
	}
	return &raw
}

func identitySession(user account.User) string {
	if sid := strings.TrimSpace(user.SessionID); sid != "" {
		return sid
	}
	if user.IsServiceCredential() {
		return strings.TrimSpace(user.CredentialID)
	}
	return ""
}

func recordAudit(operations operationlog.Recorder, user account.User, event, recordID string, detail *string, now time.Time, ctx context.Context) {
	if operations == nil {
		return
	}
	op := operationlog.Operation{
		ID:            newOperationID(),
		Event:         event,
		ActorID:       user.ID,
		ActorName:     user.Name,
		Detail:        detail,
		SessionID:     identitySession(user),
		CreatedAt:     now,
	}
	if recordID != "" {
		op.RecordID = &recordID
	}
	if ctx != nil {
		op.CorrelationID = requestid.FromContext(ctx)
	}
	if err := operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}
