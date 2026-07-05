package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

func DiscordSend(icon string, color int, withErrorInfo bool) error {
	cfg, err := loadCfg()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := validateDiscordCfg(cfg); err != nil {
		return err
	}

	env, err := loadEnv()
	if err != nil {
		return fmt.Errorf("load env: %w", err)
	}

	webhookID, err := parseDiscordWebhookID(cfg.DiscordCfg.DiscordWebhookID)
	if err != nil {
		return err
	}

	fields := discordEmbedFieldsFromNotification(buildNotificationFields(env, withErrorInfo))
	discordFields := make([]discord.EmbedField, len(fields))
	for i, field := range fields {
		discordFields[i] = discord.EmbedField{
			Name:  field.name,
			Value: field.value,
		}
	}

	client := webhook.New(webhookID, cfg.DiscordCfg.DiscordWebhookToken)
	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	title := truncateRunes(strings.TrimSpace(icon+env.Name), discordTitleMaxRunes)
	if title == "" {
		title = "Recording Notification"
	}

	if _, err := client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetEmbeds(
				discord.Embed{
					Title:  title,
					Color:  color,
					Fields: discordFields,
				},
			).Build(),
		rest.CreateWebhookMessageParams{},
		rest.WithContext(ctx),
	); err != nil {
		return fmt.Errorf("post discord message: %w", err)
	}

	return nil
}
