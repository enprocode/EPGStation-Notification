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

	for batchStart := 0; batchStart < len(batches); batchStart += discordMaxEmbedsPerMessage {
		batchEnd := batchStart + discordMaxEmbedsPerMessage
		if batchEnd > len(batches) {
			batchEnd = len(batches)
		}

		embeds := make([]discord.Embed, batchEnd-batchStart)
		for i, fields := range batches[batchStart:batchEnd] {
			embedTitle := ""
			if batchStart+i == 0 {
				embedTitle = title
			}

			discordFields := make([]discord.EmbedField, len(fields))
			for j, field := range fields {
				discordFields[j] = discord.EmbedField{
					Name:  field.name,
					Value: field.value,
				}
			}

			embeds[i] = discord.Embed{
				Title:  embedTitle,
				Color:  color,
				Fields: discordFields,
			}
		}

		if _, err := client.CreateMessage(
			discord.NewWebhookMessageCreate().WithEmbeds(embeds...),
			rest.CreateWebhookMessageParams{},
			rest.WithCtx(ctx),
		); err != nil {
			return fmt.Errorf("post discord message: %w", err)
		}
	}

	return nil
}
