package telegram

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

type telegramFailingTxRunner struct {
	err error
}

func (r telegramFailingTxRunner) Run(context.Context, func(kernel.Tx) error) error {
	return r.err
}

func TestRuntimeManager_UpdateSettingsPatchSerializesComplementaryFields(t *testing.T) {
	rm := newTestRuntimeManager(t, "initial-token", "initial-secret", nil)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		token := "updated-token"
		errs <- rm.UpdateSettingsPatch(context.Background(), &token, nil, nil, nil)
	}()
	go func() {
		defer wg.Done()
		<-start
		mode := TelegramModeWebhook
		origin := "https://console.example"
		errs <- rm.UpdateSettingsPatch(context.Background(), nil, nil, &mode, &origin)
	}()
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent UpdateSettingsPatch: %v", err)
		}
	}

	status := rm.Status()
	if rm.GetToken() != "updated-token" || rm.GetSecret() != "initial-secret" ||
		status.Mode != TelegramModeWebhook || status.WebhookPublicBaseURL != "https://console.example" {
		t.Fatalf("complementary concurrent patches lost state: token=%q secret=%q status=%+v", rm.GetToken(), rm.GetSecret(), status)
	}
}

func TestRuntimeManager_PersistenceFailureDoesNotPublishMemory(t *testing.T) {
	rm := newTestRuntimeManager(t, "initial-token", "initial-secret", nil)
	before := rm.Status()
	rm.runner = telegramFailingTxRunner{err: errors.New("database unavailable")}
	newToken := "must-not-publish"
	if err := rm.UpdateSettingsPatch(context.Background(), &newToken, nil, nil, nil); err == nil {
		t.Fatal("UpdateSettingsPatch unexpectedly succeeded with a failing persistence runner")
	}
	after := rm.Status()
	if rm.GetToken() != "initial-token" || rm.GetSecret() != "initial-secret" ||
		after.Mode != before.Mode || after.WebhookPublicBaseURL != before.WebhookPublicBaseURL {
		t.Fatalf("memory state changed after persistence failure: before=%+v after=%+v token=%q secret=%q", before, after, rm.GetToken(), rm.GetSecret())
	}
}
