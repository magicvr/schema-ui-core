package kernel

import (
	"context"
	"testing"
)

// Contract tests for the outbound-mail send port frozen by workspace-017
// GOAL-002 D-001 (R1): single bare-address recipient required, Subject and
// TextBody pass through without port-level policy.

func TestMailMessageValidate(t *testing.T) {
	tests := []struct {
		name    string
		to      string
		wantErr bool
	}{
		{"plain address", "user@example.com", false},
		{"empty recipient", "", true},
		{"whitespace-only recipient", "   ", true},
		{"missing local part", "@example.com", true},
		{"missing domain", "user@", true},
		{"display-name form rejected", "Alice <alice@example.com>", true},
		{"two addresses rejected", "a@example.com,b@example.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MailMessage{To: tt.to, Subject: "s", TextBody: "b"}
			err := msg.Validate()
			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Fatalf("MailMessage{To:%q}.Validate() error = %v, wantErr %v", tt.to, err, tt.wantErr)
			}
		})
	}
}

func TestMailMessageValidatePassesBodyThrough(t *testing.T) {
	// Port-level policy covers addressing only: empty subject/body must not
	// fail validation (product policy belongs to consumers).
	msg := MailMessage{To: "user@example.com"}
	if err := msg.Validate(); err != nil {
		t.Fatalf("Validate() with empty subject/body = %v, want nil", err)
	}
}

func TestMailSenderPortSurface(t *testing.T) {
	// Compile-time guard: the port stays a one-method synchronous contract.
	var _ MailSender = mailSenderFunc(func(ctx context.Context, msg MailMessage) error { return nil })
}

type mailSenderFunc func(ctx context.Context, msg MailMessage) error

func (f mailSenderFunc) Send(ctx context.Context, msg MailMessage) error { return f(ctx, msg) }
