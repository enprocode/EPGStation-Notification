package cmd

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

func DiscordSend(icon string, color int, withErrorInfo bool) error {
	cfg, err := loadCfg()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	env, err := loadEnv()
	if err != nil {
		return fmt.Errorf("load env: %w", err)
	}

	webhookID, err := parseDiscordWebhookID(cfg.DiscordCfg.DiscordWebhookID)
	if err != nil {
		return err
	}

	fields := buildNotificationFields(env, withErrorInfo)
	discordFields := make([]discord.EmbedField, len(fields))
	for i, field := range fields {
		discordFields[i] = discord.EmbedField{
			Name:  field.name,
			Value: field.value,
		}
	}

	client := webhook.New(webhookID, cfg.DiscordCfg.DiscordWebhookToken)
	defer client.Close(context.TODO())

	if _, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
		SetEmbeds(
			discord.Embed{
				Title:  icon + env.Name,
				Color:  color,
				Fields: discordFields,
			},
		).Build(),
	); err != nil {
		return fmt.Errorf("post discord message: %w", err)
	}

	return nil
}
