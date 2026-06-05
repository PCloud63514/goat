package slack

type SlackConfig struct {
	Mode       string `properties:"slack.mode" default:""`
	WebhookUrl string `properties:"slack.bot.url"`
	BotToken   string `properties:"slack.bot.token"`
	ChannelID  string `properties:"slack.channel.id"`
	BotName    string `properties:"slack.bot.name"`
	InfoEmoji  string `properties:"slack.info.emoji"`
	ErrorEmoji string `properties:"slack.error.emoji"`
}

type Config = SlackConfig
