package telegram

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/cache"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

type chatMemberGetterFunc func(context.Context, int64, int64) (ChatMember, error)

func (f chatMemberGetterFunc) GetChatMember(ctx context.Context, chatID, userID int64) (ChatMember, error) {
	return f(ctx, chatID, userID)
}

func boolPointer(value bool) *bool { return &value }

func newCapabilityTestService(t *testing.T, getter ChatMemberGetter) *CapabilityService {
	t.Helper()
	cachePort, err := cache.NewMemory(32)
	if err != nil {
		t.Fatal(err)
	}
	return NewCapabilityService(getter, cachePort)
}

func TestCapabilityServiceMapsMemberStatesAndCachesAllowAndDeny(t *testing.T) {
	tests := []struct {
		name   string
		member ChatMember
		want   bool
	}{
		{name: "creator", member: ChatMember{Status: "creator"}, want: true},
		{name: "member", member: ChatMember{Status: "member"}, want: true},
		{name: "administrator", member: ChatMember{Status: "administrator"}, want: true},
		{name: "administrator explicit send deny", member: ChatMember{Status: "administrator", CanSendMessages: boolPointer(false)}, want: false},
		{name: "administrator explicit post deny", member: ChatMember{Status: "administrator", CanPostMessages: boolPointer(false)}, want: false},
		{name: "restricted explicit send allow", member: ChatMember{Status: "restricted", CanSendMessages: boolPointer(true)}, want: true},
		{name: "restricted missing send permission", member: ChatMember{Status: "restricted"}, want: false},
		{name: "restricted post deny", member: ChatMember{Status: "restricted", CanSendMessages: boolPointer(true), CanPostMessages: boolPointer(false)}, want: false},
		{name: "left", member: ChatMember{Status: "left"}, want: false},
		{name: "kicked", member: ChatMember{Status: "kicked"}, want: false},
		{name: "unknown", member: ChatMember{Status: "new-status"}, want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			service := newCapabilityTestService(t, chatMemberGetterFunc(func(_ context.Context, chatID, userID int64) (ChatMember, error) {
				calls++
				if chatID != -1001 || userID != 42 {
					t.Fatalf("getChatMember ids = chat=%d user=%d, want -1001/42", chatID, userID)
				}
				return test.member, nil
			}))

			got, err := service.Check(context.Background(), 42, -1001, false)
			if err != nil || got != test.want {
				t.Fatalf("first Check = %v, %v; want %v, nil", got, err, test.want)
			}
			got, err = service.Check(context.Background(), 42, -1001, false)
			if err != nil || got != test.want || calls != 1 {
				t.Fatalf("cached Check = %v, %v calls=%d; want %v, nil, 1", got, err, calls, test.want)
			}
		})
	}
}

func TestCapabilityServiceForceRefreshAndForbiddenAreCached(t *testing.T) {
	var mu sync.Mutex
	member := ChatMember{Status: "member"}
	calls := 0
	service := newCapabilityTestService(t, chatMemberGetterFunc(func(_ context.Context, chatID, userID int64) (ChatMember, error) {
		if chatID != 8001 || userID != 42 {
			t.Fatalf("getChatMember ids = chat=%d user=%d, want 8001/42", chatID, userID)
		}
		mu.Lock()
		defer mu.Unlock()
		calls++
		return member, nil
	}))

	allowed, err := service.Check(context.Background(), 42, 8001, false)
	if err != nil || !allowed {
		t.Fatalf("initial capability = %v, %v; want true, nil", allowed, err)
	}
	member = ChatMember{Status: "left"}
	allowed, err = service.Check(context.Background(), 42, 8001, false)
	if err != nil || !allowed {
		t.Fatalf("non-forced cached capability = %v, %v; want cached true, nil", allowed, err)
	}
	allowed, err = service.Check(context.Background(), 42, 8001, true)
	if err != nil || allowed {
		t.Fatalf("forced capability = %v, %v; want false, nil", allowed, err)
	}
	allowed, err = service.Check(context.Background(), 42, 8001, false)
	if err != nil || allowed {
		t.Fatalf("post-force cached capability = %v, %v; want cached false, nil", allowed, err)
	}
	mu.Lock()
	if calls != 2 {
		t.Fatalf("getChatMember calls = %d, want 2", calls)
	}
	mu.Unlock()
}

func TestCapabilityServiceCacheIsScopedByBotAndChat(t *testing.T) {
	calls := make(map[[2]int64]int)
	service := newCapabilityTestService(t, chatMemberGetterFunc(func(_ context.Context, chatID, userID int64) (ChatMember, error) {
		key := [2]int64{userID, chatID}
		calls[key]++
		if userID == 42 && chatID == 8001 {
			return ChatMember{Status: "member"}, nil
		}
		return ChatMember{Status: "left"}, nil
	}))

	checks := []struct {
		botID  int64
		chatID int64
		want   bool
	}{
		{botID: 42, chatID: 8001, want: true},
		{botID: 43, chatID: 8001, want: false},
		{botID: 42, chatID: 8002, want: false},
	}
	for _, check := range checks {
		got, err := service.Check(context.Background(), check.botID, check.chatID, false)
		if err != nil || got != check.want {
			t.Fatalf("capability bot=%d chat=%d = %v, %v; want %v, nil", check.botID, check.chatID, got, err, check.want)
		}
	}
	if got, err := service.Check(context.Background(), 42, 8001, false); err != nil || !got {
		t.Fatalf("cached scoped capability = %v, %v; want true, nil", got, err)
	}
	if calls[[2]int64{42, 8001}] != 1 || calls[[2]int64{43, 8001}] != 1 || calls[[2]int64{42, 8002}] != 1 {
		t.Fatalf("getChatMember calls by bot/chat = %v, want one per distinct key", calls)
	}
}

func TestCapabilityServiceSingleFlightAndForbiddenMapping(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var startOnce sync.Once
	var mu sync.Mutex
	calls := 0
	service := newCapabilityTestService(t, chatMemberGetterFunc(func(context.Context, int64, int64) (ChatMember, error) {
		mu.Lock()
		calls++
		mu.Unlock()
		startOnce.Do(func() { close(started) })
		<-release
		return ChatMember{Status: "member"}, nil
	}))

	type result struct {
		allowed bool
		err     error
	}
	results := make(chan result, 2)
	go func() {
		allowed, err := service.Check(context.Background(), 42, 8001, false)
		results <- result{allowed: allowed, err: err}
	}()
	<-started
	go func() {
		allowed, err := service.Check(context.Background(), 42, 8001, true)
		results <- result{allowed: allowed, err: err}
	}()
	deadline := time.Now().Add(time.Second)
	for {
		service.mu.Lock()
		flight := service.flights[capabilityCacheKey(42, 8001)]
		joined := flight != nil && flight.waiters == 1
		service.mu.Unlock()
		if joined {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("second Check did not join the in-flight capability probe")
		}
		time.Sleep(time.Millisecond)
	}
	close(release)
	for range 2 {
		got := <-results
		if got.err != nil || !got.allowed {
			t.Fatalf("single-flight result = %v, %v; want true, nil", got.allowed, got.err)
		}
	}
	mu.Lock()
	if calls != 1 {
		t.Fatalf("single-flight getChatMember calls = %d, want 1", calls)
	}
	mu.Unlock()

	forbidden := newCapabilityTestService(t, chatMemberGetterFunc(func(context.Context, int64, int64) (ChatMember, error) {
		return ChatMember{}, &TelegramAPIError{Method: "getChatMember", HTTPStatus: 403, ErrorCode: 403}
	}))
	allowed, err := forbidden.Check(context.Background(), 42, 8002, false)
	if err != nil || allowed {
		t.Fatalf("forbidden capability = %v, %v; want false, nil", allowed, err)
	}
	allowed, err = forbidden.Check(context.Background(), 42, 8002, false)
	if err != nil || allowed {
		t.Fatalf("cached forbidden capability = %v, %v; want false, nil", allowed, err)
	}
}

type telegramSenderFunc func(context.Context, kernel.TelegramMessage) error

func (f telegramSenderFunc) Send(ctx context.Context, message kernel.TelegramMessage) error {
	return f(ctx, message)
}

func TestCapabilityInvalidatingSenderDeletesExactBotChatEntry(t *testing.T) {
	member := ChatMember{Status: "member"}
	calls := 0
	service := newCapabilityTestService(t, chatMemberGetterFunc(func(context.Context, int64, int64) (ChatMember, error) {
		calls++
		return member, nil
	}))
	if allowed, err := service.Check(context.Background(), 42, 8001, false); err != nil || !allowed {
		t.Fatalf("seed capability = %v, %v; want true, nil", allowed, err)
	}
	member = ChatMember{Status: "left"}
	sender := NewCapabilityInvalidatingSender(
		telegramSenderFunc(func(context.Context, kernel.TelegramMessage) error {
			return &TelegramAPIError{Method: "sendMessage", HTTPStatus: 200, ErrorCode: 403}
		}),
		service,
		func() int64 { return 42 },
	)
	if err := sender.Send(context.Background(), kernel.TelegramMessage{ChatID: "8001", Text: "hello"}); err == nil {
		t.Fatal("403 sender unexpectedly succeeded")
	}
	if allowed, err := service.Check(context.Background(), 42, 8001, false); err != nil || allowed {
		t.Fatalf("post-403 capability = %v, %v; want fresh false, nil", allowed, err)
	}
	if calls != 2 {
		t.Fatalf("getChatMember calls after exact invalidation = %d, want 2", calls)
	}

	if allowed, err := service.Check(context.Background(), 42, 8001, false); err != nil || allowed {
		t.Fatalf("denial cache = %v, %v; want false, nil", allowed, err)
	}
	if err := (NewCapabilityInvalidatingSender(
		telegramSenderFunc(func(context.Context, kernel.TelegramMessage) error { return errors.New("temporary failure") }),
		service,
		func() int64 { return 42 },
	)).Send(context.Background(), kernel.TelegramMessage{ChatID: "8001", Text: "hello"}); err == nil {
		t.Fatal("non-403 sender unexpectedly succeeded")
	}
	if calls != 2 {
		t.Fatalf("non-403 sender changed cached capability; calls=%d, want 2", calls)
	}
}
