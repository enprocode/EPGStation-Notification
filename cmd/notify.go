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
	slackFieldMaxRunes        = 2000
	slackTitleMaxRunes        = 256
	discordFieldNameMaxRunes  = 256
	discordFieldMaxRunes      = 1024
	discordTitleMaxRunes      = 256
	discordMaxEmbedFields     = 25
	discordMaxEmbedTotalRunes = 6000
	discordMaxEmbedsPerMessage = 10
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

func formatChannel(name, channelType string) string {
	name = strings.TrimSpace(name)
	channelType = strings.TrimSpace(channelType)

	switch {
	case name != "" && channelType != "":
		return name + "/" + channelType
	case name != "":
		return name
	case channelType != "":
		return channelType
	default:
		return ""
	}
}

func buildNotificationFields(env cmdEnv, withErrorInfo bool) []notificationField {
	start, end := formatProgramTime(env.StartAt, env.EndAt)

	fields := make([]notificationField, 0, 8)
	fields = appendFieldIfPresent(fields, "Channel", formatChannel(env.ChannelName, env.ChannelType))
	fields = appendFieldIfPresent(fields, "Time", start+" ~ "+end)
	fields = appendFieldIfPresent(fields, "Description", env.Description)
	fields = appendFieldIfPresent(fields, "Extended", env.Extended)
	fields = appendFieldIfPresent(fields, "RecPath", env.RecPath)

	if withErrorInfo {
		fields = appendFieldIfPresent(fields, "Error Count", env.ErrorCnt)
		fields = appendFieldIfPresent(fields, "Drop Count", env.DropCnt)
		fields = appendFieldIfPresent(fields, "LogPath", env.LogPath)
	}

	return fields
}

func appendFieldIfPresent(fields []notificationField, name, value string) []notificationField {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" || value == "" {
		return fields
	}
	return append(fields, notificationField{name: name, value: value})
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

func splitDiscordField(name, value string) []notificationField {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" || value == "" {
		return nil
	}

	runes := []rune(value)
	if len(runes) <= discordFieldMaxRunes {
		return []notificationField{{name: name, value: value}}
	}

	totalParts := (len(runes) + discordFieldMaxRunes - 1) / discordFieldMaxRunes
	parts := make([]notificationField, 0, totalParts)
	for i, start := 0, 0; start < len(runes); i++ {
		end := start + discordFieldMaxRunes
		if end > len(runes) {
			end = len(runes)
		}

		partName := name
		if totalParts > 1 {
			partName = fmt.Sprintf("%s (%d/%d)", name, i+1, totalParts)
		}

		parts = append(parts, notificationField{
			name:  truncateRunes(partName, discordFieldNameMaxRunes),
			value: string(runes[start:end]),
		})
		start = end
	}

	return parts
}

func discordEmbedTextRunes(title string, fields []notificationField) int {
	total := len([]rune(title))
	for _, field := range fields {
		total += len([]rune(field.name)) + len([]rune(field.value))
	}
	return total
}

func discordEmbedBatchesFromNotification(title string, fields []notificationField) [][]notificationField {
	var batches [][]notificationField
	current := make([]notificationField, 0, len(fields))
	firstBatch := true

	flush := func() {
		if len(current) == 0 {
			return
		}
		batches = append(batches, current)
		current = make([]notificationField, 0, cap(current))
		firstBatch = false
	}

	embedTitle := func() string {
		if firstBatch {
			return title
		}
		return ""
	}

	fits := func(part notificationField) bool {
		if len(current) >= discordMaxEmbedFields {
			return false
		}
		candidate := append(current, part)
		return discordEmbedTextRunes(embedTitle(), candidate) <= discordMaxEmbedTotalRunes
	}

	addPart := func(part notificationField) {
		for {
			if fits(part) {
				current = append(current, part)
				return
			}

			if len(current) > 0 {
				flush()
				continue
			}

			nameRunes := len([]rune(part.name))
			maxValueRunes := discordMaxEmbedTotalRunes - discordEmbedTextRunes(embedTitle(), nil) - nameRunes
			if maxValueRunes <= 0 {
				maxValueRunes = 1
			}
			part.value = truncateRunes(part.value, maxValueRunes)
			current = append(current, part)
			return
		}
	}

	for _, field := range fields {
		for _, part := range splitDiscordField(field.name, field.value) {
			addPart(part)
		}
	}

	flush()
	return batches
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
