package slack

import (
	"errors"
	"testing"
)

func TestNewSenderModeWebhook(t *testing.T) {
	s := New(&SlackConfig{
		Mode:       "webhook",
		BotToken:   "xoxb-test",
		WebhookUrl: "https://example.com/hook",
	})

	if !s.initialized {
		t.Fatal("expected slack to be initialized")
	}

	if s.Mode() != "webhook" {
		t.Fatalf("expected webhook mode, got %q", s.Mode())
	}

	if _, ok := s.senders["webhook"].(webhookSender); !ok {
		t.Fatalf("expected webhook sender, got %T", s.senders["webhook"])
	}
}

func TestNewSenderModeChat(t *testing.T) {
	s := New(&SlackConfig{
		Mode:     "chat",
		BotToken: "xoxb-test",
	})

	if !s.initialized {
		t.Fatal("expected slack to be initialized")
	}

	if s.Mode() != "chat" {
		t.Fatalf("expected chat mode, got %q", s.Mode())
	}

	if _, ok := s.senders["chat"].(chatSender); !ok {
		t.Fatalf("expected chat sender, got %T", s.senders["chat"])
	}
}

func TestNewSenderDefaultPrefersChat(t *testing.T) {
	s := New(&SlackConfig{
		BotToken:   "xoxb-test",
		WebhookUrl: "https://example.com/hook",
	})

	if !s.initialized {
		t.Fatal("expected slack to be initialized")
	}

	if s.Mode() != "chat" {
		t.Fatalf("expected chat mode, got %q", s.Mode())
	}
}

func TestWithFieldsKeepsSender(t *testing.T) {
	base := New(&SlackConfig{
		Mode:       "webhook",
		WebhookUrl: "https://example.com/hook",
	})

	withFields := base.WithFields(AttachmentField{Key: "k", Value: "v"})

	if len(withFields.fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(withFields.fields))
	}

	if withFields.Mode() != "webhook" {
		t.Fatalf("expected webhook mode, got %q", withFields.Mode())
	}

	if _, ok := withFields.senders["webhook"].(webhookSender); !ok {
		t.Fatalf("expected webhook sender, got %T", withFields.senders["webhook"])
	}
}

func TestSetModeSwitchesActiveSender(t *testing.T) {
	s := New(&SlackConfig{
		BotToken:   "xoxb-test",
		WebhookUrl: "https://example.com/hook",
	})

	if err := (&s).SetMode("webhook"); err != nil {
		t.Fatalf("expected mode switch to succeed: %v", err)
	}

	if s.Mode() != "webhook" {
		t.Fatalf("expected webhook mode, got %q", s.Mode())
	}
}

func TestSetModeRejectsUnavailableSender(t *testing.T) {
	s := New(&SlackConfig{
		BotToken: "xoxb-test",
	})

	if err := (&s).SetMode("webhook"); err == nil {
		t.Fatal("expected mode switch to fail")
	} else if !errors.Is(err, ErrModeUnavailable) {
		t.Fatalf("expected ErrModeUnavailable, got %v", err)
	}
}

func TestNotifyReturnsNotConfiguredWhenModeMissing(t *testing.T) {
	s := New(&SlackConfig{})

	if err := s.Info("hello"); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestExportedInfoReturnsNotConfigured(t *testing.T) {
	slack = New(nil)

	if err := Info("hello"); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestExportedSetModeReturnsNotConfigured(t *testing.T) {
	slack = New(nil)

	if err := SetMode("chat"); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}
