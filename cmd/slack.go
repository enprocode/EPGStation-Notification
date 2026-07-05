package cmd

import (
	"fmt"
	"net/http"
	"strings"

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
	slackFields := make([]slack.AttachmentField, 0, len(fields))
	for _, field := range fields {
		name := strings.TrimSpace(field.name)
		value := strings.TrimSpace(field.value)
		if name == "" || value == "" {
			continue
		}
		slackFields = append(slackFields, slack.AttachmentField{
			Title: truncateRunes(name, slackFieldMaxRunes),
			Value: truncateRunes(value, slackFieldMaxRunes),
			Short: false,
		})
	}

	httpClient := &http.Client{Timeout: requestTimeout}
	api := slack.New(cfg.SlackCfg.SlackToken, slack.OptionHTTPClient(httpClient))
	title := truncateRunes(strings.TrimSpace(icon+env.Name), slackTitleMaxRunes)
	if title == "" {
		title = "Recording Notification"
	}
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
