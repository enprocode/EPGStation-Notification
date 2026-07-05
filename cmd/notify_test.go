package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFormatProgramTime(t *testing.T) {
	startAt := 1728000000000
	endAt := 1728003600000
	start, end := formatProgramTime(startAt, endAt)

	wantStart := time.Unix(int64(startAt/1000), 0).Format("2006-01-02 15:04 MST")
	wantEnd := time.Unix(int64(endAt/1000), 0).Format("2006-01-02 15:04 MST")
	if start != wantStart || end != wantEnd {
		t.Fatalf("got start=%q end=%q, want start=%q end=%q", start, end, wantStart, wantEnd)
	}

	zeroStart, zeroEnd := formatProgramTime(0, 0)
	if zeroStart != "-" || zeroEnd != "-" {
		t.Fatalf("expected '-' for zero timestamps, got start=%q end=%q", zeroStart, zeroEnd)
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := truncateRunes("hello", 10); got != "hello" {
		t.Fatalf("unexpected truncate result: %q", got)
	}

	long := strings.Repeat("あ", 10)
	if got := truncateRunes(long, 7); got != strings.Repeat("あ", 4)+"..." {
		t.Fatalf("unexpected truncate result: %q", got)
	}

	if got := truncateRunes("hello", 0); got != "" {
		t.Fatalf("expected empty string for non-positive max, got %q", got)
	}
	if got := truncateRunes("hello", 3); got != "hel" {
		t.Fatalf("expected hard cut without ellipsis, got %q", got)
	}
}

func TestFormatChannel(t *testing.T) {
	if got := formatChannel("NHK", "GR"); got != "NHK/GR" {
		t.Fatalf("unexpected channel: %q", got)
	}
	if got := formatChannel("  ", ""); got != "" {
		t.Fatalf("expected empty channel, got %q", got)
	}
	if got := formatChannel("NHK", ""); got != "NHK" {
		t.Fatalf("expected name only, got %q", got)
	}
	if got := formatChannel("", "GR"); got != "GR" {
		t.Fatalf("expected type only, got %q", got)
	}
}

func TestBuildNotificationFieldsSkipsWhitespaceDescription(t *testing.T) {
	fields := buildNotificationFields(cmdEnv{
		ChannelName: "NHK",
		ChannelType: "GR",
		StartAt:     1728000000000,
		EndAt:       1728003600000,
		Description: "   ",
	}, false)

	for _, field := range fields {
		if field.name == "Description" {
			t.Fatal("expected whitespace-only description to be skipped")
		}
	}
}

func TestSplitDiscordField(t *testing.T) {
	longValue := strings.Repeat("a", discordFieldMaxRunes+100)
	parts := splitDiscordField("Description", longValue)
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if len([]rune(parts[0].value)) != discordFieldMaxRunes {
		t.Fatalf("unexpected first part length: %d", len([]rune(parts[0].value)))
	}
	if parts[0].name != "Description (1/2)" {
		t.Fatalf("unexpected part name: %q", parts[0].name)
	}
}

func TestSplitDiscordFieldSkipsEmpty(t *testing.T) {
	if parts := splitDiscordField("Description", "   "); parts != nil {
		t.Fatalf("expected nil parts for empty value, got %v", parts)
	}
}

func TestDiscordEmbedBatchesFromNotification(t *testing.T) {
	fields := discordEmbedBatchesFromNotification("title", []notificationField{
		{name: "Description", value: "   "},
		{name: "RecPath", value: "/recordings/test.m2ts"},
		{name: "Extended", value: strings.Repeat("x", discordFieldMaxRunes+10)},
	})

	if len(fields) != 1 || len(fields[0]) != 3 {
		t.Fatalf("expected 1 batch with 3 fields, got %d batches", len(fields))
	}
	if fields[0][0].name != "RecPath" {
		t.Fatalf("expected RecPath first, got %q", fields[0][0].name)
	}
	if fields[0][1].name != "Extended (1/2)" {
		t.Fatalf("unexpected split field name: %q", fields[0][1].name)
	}
}

func TestDiscordEmbedBatchesRespectTotalRunes(t *testing.T) {
	title := strings.Repeat("t", discordTitleMaxRunes)
	longValue := strings.Repeat("a", discordFieldMaxRunes)
	fields := make([]notificationField, 0, 7)
	for i := 1; i <= 7; i++ {
		fields = append(fields, notificationField{
			name:  fmt.Sprintf("Field %d", i),
			value: longValue,
		})
	}

	batches := discordEmbedBatchesFromNotification(title, fields)
	if len(batches) < 2 {
		t.Fatalf("expected multiple embed batches, got %d", len(batches))
	}

	for i, batch := range batches {
		embedTitle := ""
		if i == 0 {
			embedTitle = title
		}
		if got := discordEmbedTextRunes(embedTitle, batch); got > discordMaxEmbedTotalRunes {
			t.Fatalf("batch %d exceeds embed rune limit: %d", i, got)
		}
		if len(batch) > discordMaxEmbedFields {
			t.Fatalf("batch %d exceeds field limit: %d", i, len(batch))
		}
	}
}

func TestDiscordEmbedBatchesSplitAcrossManyParts(t *testing.T) {
	longValue := strings.Repeat("x", discordFieldMaxRunes*8)
	batches := discordEmbedBatchesFromNotification("Program", []notificationField{
		{name: "Description", value: longValue},
	})

	totalFields := 0
	for i, batch := range batches {
		totalFields += len(batch)
		embedTitle := "Program"
		if i > 0 {
			embedTitle = ""
		}
		if got := discordEmbedTextRunes(embedTitle, batch); got > discordMaxEmbedTotalRunes {
			t.Fatalf("batch %d exceeds embed rune limit: %d", i, got)
		}
	}
	if totalFields < 8 {
		t.Fatalf("expected at least 8 field parts across batches, got %d", totalFields)
	}
}

func TestBuildNotificationFieldsWithErrorInfo(t *testing.T) {
	fields := buildNotificationFields(cmdEnv{
		ChannelName: "NHK",
		ChannelType: "GR",
		StartAt:     1728000000000,
		EndAt:       1728003600000,
		Description: "test program",
		RecPath:     "/recordings/test.m2ts",
		ErrorCnt:    "3",
		DropCnt:     "1",
		LogPath:     "/logs/rec.log",
	}, true)

	if len(fields) != 7 {
		t.Fatalf("expected 7 fields, got %d", len(fields))
	}
	if fields[len(fields)-1].name != "LogPath" {
		t.Fatalf("expected LogPath field, got %s", fields[len(fields)-1].name)
	}
}

func TestParseDiscordWebhookID(t *testing.T) {
	id, err := parseDiscordWebhookID("1234567890123456789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.String() != "1234567890123456789" {
		t.Fatalf("unexpected id: %s", id.String())
	}

	if _, err := parseDiscordWebhookID(""); err == nil {
		t.Fatal("expected error for empty id")
	}

	if _, err := parseDiscordWebhookID("not-a-number"); err == nil {
		t.Fatal("expected error for non-numeric id")
	}
}

func TestValidateSlackCfg(t *testing.T) {
	if err := validateSlackCfg(cmdCfg{}); err == nil {
		t.Fatal("expected error for empty slack config")
	}

	tokenOnly := cmdCfg{}
	tokenOnly.SlackCfg.SlackToken = "xoxb-token"
	if err := validateSlackCfg(tokenOnly); err == nil {
		t.Fatal("expected error when slack channel is missing")
	}

	valid := tokenOnly
	valid.SlackCfg.Channel = "C1234567890"
	if err := validateSlackCfg(valid); err != nil {
		t.Fatalf("unexpected error for valid slack config: %v", err)
	}
}

func TestValidateDiscordCfg(t *testing.T) {
	if err := validateDiscordCfg(cmdCfg{}); err == nil {
		t.Fatal("expected error for empty discord config")
	}

	badID := cmdCfg{}
	badID.DiscordCfg.DiscordWebhookToken = "webhook-token"
	badID.DiscordCfg.DiscordWebhookID = "not-a-number"
	if err := validateDiscordCfg(badID); err == nil {
		t.Fatal("expected error for invalid webhook id")
	}

	valid := badID
	valid.DiscordCfg.DiscordWebhookID = "1234567890123456789"
	if err := validateDiscordCfg(valid); err != nil {
		t.Fatalf("unexpected error for valid discord config: %v", err)
	}
}
