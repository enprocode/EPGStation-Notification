package cmd

import (
	"fmt"
	"net/http"
	"time"

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
	if err := validateSlackCfg(cfg); err != nil {
		return err
	}

	fields := buildNotificationFields(env, withErrorInfo)
	slackFields := make([]slack.AttachmentField, len(fields))
	for i, field := range fields {
		slackFields[i] = slack.AttachmentField{
			Title: truncateRunes(field.name, slackFieldMaxRunes),
			Value: truncateRunes(field.value, slackFieldMaxRunes),
			Short: false,
		}
	}

	httpClient := &http.Client{Timeout: requestTimeout}
	api := slack.New(cfg.SlackCfg.SlackToken, slack.OptionHTTPClient(httpClient))
	title := truncateRunes(icon+env.Name, slackTitleMaxRunes)
	attachment := slack.Attachment{
		Fallback: title,
		Color:    color,
		Title:    title,
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
