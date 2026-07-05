package cmd

import (
	"context"
	"fmt"
	"net/http"
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

	title := truncateRunes(strings.TrimSpace(icon+env.Name), discordTitleMaxRunes)
	if title == "" {
		title = "Recording Notification"
	}

	batches := discordEmbedBatchesFromNotification(title, buildNotificationFields(env, withErrorInfo))
	if len(batches) == 0 {
		batches = [][]notificationField{{}}
	}

	httpClient := &http.Client{Timeout: requestTimeout}
	client := webhook.New(
		webhookID,
		cfg.DiscordCfg.DiscordWebhookToken,
		webhook.WithRestClientConfigOpts(rest.WithHTTPClient(httpClient)),
	)
	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Discord's 6000-character limit applies to the combined text of all embeds
	// in a single message. Since each batch is already capped at that limit, we
	// send one embed per message to stay within it.
	for i, fields := range batches {
		embedTitle := ""
		if i == 0 {
			embedTitle = title
		}

		discordFields := make([]discord.EmbedField, len(fields))
		for j, field := range fields {
			discordFields[j] = discord.EmbedField{
				Name:  field.name,
				Value: field.value,
			}
		}

		embed := discord.Embed{
			Title:  embedTitle,
			Color:  color,
			Fields: discordFields,
		}

		if _, err := client.CreateMessage(
			discord.NewWebhookMessageCreate().WithEmbeds(embed),
			rest.CreateWebhookMessageParams{},
			rest.WithCtx(ctx),
		); err != nil {
			return fmt.Errorf("post discord message: %w", err)
		}
	}

	return nil
}
