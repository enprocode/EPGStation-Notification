package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

const requestTimeout = 30 * time.Second

const (
	slackFieldMaxRunes  = 2000
	slackTitleMaxRunes  = 256
	discordFieldMaxRunes = 1024
	discordTitleMaxRunes = 256
)

type notificationField struct {
	name  string
	value string
}

func formatProgramTime(startAt, endAt int) (string, string) {
	format := func(ts int) string {
		if ts <= 0 {
			return "-"
		}
		return time.Unix(int64(ts/1000), 0).Format("2006-01-02 15:04 MST")
	}
	return format(startAt), format(endAt)
}

func buildNotificationFields(env cmdEnv, withErrorInfo bool) []notificationField {
	start, end := formatProgramTime(env.StartAt, env.EndAt)

	fields := []notificationField{
		{name: "Channel", value: env.ChannelName + "/" + env.ChannelType},
		{name: "Time", value: start + " ~ " + end},
	}

	if env.Description != "" {
		fields = append(fields, notificationField{name: "Description", value: env.Description})
	}
	if env.Extended != "" {
		fields = append(fields, notificationField{name: "Extended", value: env.Extended})
	}
	if env.RecPath != "" {
		fields = append(fields, notificationField{name: "RecPath", value: env.RecPath})
	}

	if withErrorInfo {
		if env.ErrorCnt != "" {
			fields = append(fields, notificationField{name: "Error Count", value: env.ErrorCnt})
		}
		if env.DropCnt != "" {
			fields = append(fields, notificationField{name: "Drop Count", value: env.DropCnt})
		}
		if env.LogPath != "" {
			fields = append(fields, notificationField{name: "LogPath", value: env.LogPath})
		}
	}

	return fields
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 3 {
		return string(runes[:maxRunes])
	}

	return string(runes[:maxRunes-3]) + "..."
}

func parseDiscordWebhookID(id string) (snowflake.ID, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, fmt.Errorf("discord webhook id is empty")
	}

	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid discord webhook id %q: %w", id, err)
	}

	return snowflake.ID(parsed), nil
}
