package slack

import (
	"errors"
	"fmt"
	"strings"

	slackapi "github.com/slack-go/slack"
)

type Slack struct {
	initialized bool
	mode        string
	channelID   string
	botName     string
	infoEmoji   string
	errorEmoji  string
	senders     map[string]sender
	fields      []AttachmentField
}

var (
	ErrNotConfigured   = errors.New("slack is not configured")
	ErrInvalidMode     = errors.New("invalid slack mode")
	ErrModeUnavailable = errors.New("slack mode is not configured")
)

func New(config *SlackConfig) Slack {
	if config == nil {
		return Slack{}
	}

	senders := newSenders(config)
	mode := resolveMode(config, senders)

	return Slack{
		initialized: len(senders) > 0 && mode != "",
		mode:        mode,
		channelID:   config.ChannelID,
		botName:     config.BotName,
		infoEmoji:   config.InfoEmoji,
		errorEmoji:  config.ErrorEmoji,
		senders:     senders,
	}
}

func (s Slack) Info(args ...interface{}) error {
	return s.notify(fmt.Sprint(args...), s.infoEmoji)
}

func (s Slack) Infof(format string, args ...interface{}) error {
	return s.notify(fmt.Sprintf(format, args...), s.infoEmoji)
}

func (s Slack) Error(args ...interface{}) error {
	return s.notify(fmt.Sprint(args...), s.errorEmoji)
}

func (s Slack) Errorf(format string, args ...interface{}) error {
	return s.notify(fmt.Sprintf(format, args...), s.errorEmoji)
}

func (s Slack) WithFields(fields ...AttachmentField) Slack {
	combinedFields := make([]AttachmentField, 0, len(s.fields)+len(fields))
	combinedFields = append(combinedFields, s.fields...)
	combinedFields = append(combinedFields, fields...)

	return Slack{
		initialized: s.initialized,
		mode:        s.mode,
		channelID:   s.channelID,
		botName:     s.botName,
		infoEmoji:   s.infoEmoji,
		errorEmoji:  s.errorEmoji,
		senders:     s.senders,
		fields:      combinedFields,
	}
}

func (s *Slack) SetMode(mode string) error {
	normalized := normalizeMode(mode)
	if normalized == "" {
		return fmt.Errorf("%w: %q", ErrInvalidMode, mode)
	}
	if _, ok := s.senders[normalized]; !ok {
		return fmt.Errorf("%w: %q", ErrModeUnavailable, normalized)
	}

	s.mode = normalized
	s.initialized = true
	return nil
}

func (s Slack) Mode() string {
	return s.mode
}

type AttachmentField struct {
	Key   string
	Value interface{}
}

func (s Slack) notify(message string, emoji string) error {
	sender := s.senders[s.mode]
	if sender == nil {
		if s.mode == "" {
			return ErrNotConfigured
		}
		return fmt.Errorf("%w: %q", ErrModeUnavailable, s.mode)
	}

	err := sender.Notify(notification{
		channelID: s.channelID,
		botName:   s.botName,
		emoji:     emoji,
		message:   message,
		fields:    s.fields,
	})
	if err != nil {
		return fmt.Errorf("slack notify failed in %s mode: %w", s.mode, err)
	}

	return nil
}

func stringify(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

type notification struct {
	channelID string
	botName   string
	emoji     string
	message   string
	fields    []AttachmentField
}

type sender interface {
	Notify(notification) error
}

func newSenders(config *SlackConfig) map[string]sender {
	senders := map[string]sender{}

	if config.BotToken != "" {
		senders["chat"] = chatSender{api: slackapi.New(config.BotToken)}
	}
	if config.WebhookUrl != "" {
		senders["webhook"] = webhookSender{url: config.WebhookUrl}
	}

	return senders
}

func resolveMode(config *SlackConfig, senders map[string]sender) string {
	mode := normalizeMode(config.Mode)
	if mode != "" {
		if _, ok := senders[mode]; ok {
			return mode
		}
		return ""
	}

	if _, ok := senders["chat"]; ok {
		return "chat"
	}
	if _, ok := senders["webhook"]; ok {
		return "webhook"
	}

	return ""
}

func normalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "chat":
		return "chat"
	case "webhook":
		return "webhook"
	default:
		return ""
	}
}

type webhookSender struct {
	url string
}

func (s webhookSender) Notify(n notification) error {
	payload := &slackapi.WebhookMessage{
		Channel:     n.channelID,
		Username:    n.botName,
		IconEmoji:   n.emoji,
		IconURL:     n.emoji,
		Text:        n.message,
		Attachments: []slackapi.Attachment{},
	}
	if len(n.fields) > 0 {
		attachment := slackapi.Attachment{}
		for i := range n.fields {
			attachment.Fields = append(attachment.Fields, slackapi.AttachmentField{
				Title: n.fields[i].Key,
				Value: stringify(n.fields[i].Value),
				Short: true,
			})
		}
		payload.Attachments = append(payload.Attachments, attachment)
	}

	return slackapi.PostWebhook(s.url, payload)
}

type chatSender struct {
	api *slackapi.Client
}

func (s chatSender) Notify(n notification) error {
	options := []slackapi.MsgOption{
		slackapi.MsgOptionText(n.message, false),
	}

	if n.botName != "" {
		options = append(options, slackapi.MsgOptionUsername(n.botName))
	}
	if n.emoji != "" {
		options = append(options, slackapi.MsgOptionIconEmoji(n.emoji))
	}
	if len(n.fields) > 0 {
		attachment := slackapi.Attachment{}
		for i := range n.fields {
			attachment.Fields = append(attachment.Fields, slackapi.AttachmentField{
				Title: n.fields[i].Key,
				Value: stringify(n.fields[i].Value),
				Short: true,
			})
		}
		options = append(options, slackapi.MsgOptionAttachments(attachment))
	}

	_, _, err := s.api.PostMessage(n.channelID, options...)
	return err
}
