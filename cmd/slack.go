package cmd

import (
	"fmt"

	"github.com/slack-go/slack"
)

func Slack(icon, color string, withErrorInfo bool) error {
	env, err := loadEnv()
	if err != nil {
		return fmt.Errorf("load env: %w", err)
	}

	cfg, err := loadCfg()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fields := buildNotificationFields(env, withErrorInfo)
	slackFields := make([]slack.AttachmentField, len(fields))
	for i, field := range fields {
		slackFields[i] = slack.AttachmentField{
			Title: field.name,
			Value: field.value,
			Short: false,
		}
	}

	api := slack.New(cfg.SlackCfg.SlackToken)
	attachment := slack.Attachment{
		Fallback: icon + env.Name,
		Color:    color,
		Title:    icon + env.Name,
		Fields:   slackFields,
	}

	_, _, err = api.PostMessage(
		cfg.SlackCfg.Channel,
		slack.MsgOptionAsUser(false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return fmt.Errorf("post slack message: %w", err)
	}

	return nil
}
