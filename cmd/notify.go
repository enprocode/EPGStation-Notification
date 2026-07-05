package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type notificationField struct {
	name  string
	value string
}

func formatProgramTime(startAt, endAt int) (string, string) {
	start := time.Unix(int64(startAt/1000), 0).Format("2006-01-02 15:04")
	end := time.Unix(int64(endAt/1000), 0).Format("2006-01-02 15:04")
	return start, end
}

func buildNotificationFields(env cmdEnv, withErrorInfo bool) []notificationField {
	start, end := formatProgramTime(env.StartAt, env.EndAt)

	fields := []notificationField{
		{name: "Channel", value: env.ChannelName + "/" + env.ChannelType},
		{name: "Time", value: start + " ~ " + end},
		{name: "Description", value: env.Description},
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
