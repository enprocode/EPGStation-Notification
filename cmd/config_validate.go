package cmd

import (
	"fmt"
	"strings"
)

func validateSlackCfg(cfg cmdCfg) error {
	if strings.TrimSpace(cfg.SlackCfg.SlackToken) == "" {
		return fmt.Errorf("slack-token is not configured in config.yml")
	}
	if strings.TrimSpace(cfg.SlackCfg.Channel) == "" {
		return fmt.Errorf("slack channel is not configured in config.yml")
	}
	return nil
}

func validateDiscordCfg(cfg cmdCfg) error {
	if strings.TrimSpace(cfg.DiscordCfg.DiscordWebhookToken) == "" {
		return fmt.Errorf("discord-webhook-token is not configured in config.yml")
	}
	if _, err := parseDiscordWebhookID(cfg.DiscordCfg.DiscordWebhookID); err != nil {
		return err
	}
	return nil
}
