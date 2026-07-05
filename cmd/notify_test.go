package cmd

import (
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
}

func TestFormatChannel(t *testing.T) {
	if got := formatChannel("NHK", "GR"); got != "NHK/GR" {
		t.Fatalf("unexpected channel: %q", got)
	}
	if got := formatChannel("  ", ""); got != "" {
		t.Fatalf("expected empty channel, got %q", got)
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

func TestDiscordEmbedFieldsFromNotification(t *testing.T) {
	fields := discordEmbedFieldsFromNotification([]notificationField{
		{name: "Description", value: "   "},
		{name: "RecPath", value: "/recordings/test.m2ts"},
		{name: "Extended", value: strings.Repeat("x", discordFieldMaxRunes+10)},
	})

	if len(fields) != 3 {
		t.Fatalf("expected 3 sanitized fields, got %d", len(fields))
	}
	if fields[0].name != "RecPath" {
		t.Fatalf("expected RecPath first, got %q", fields[0].name)
	}
	if fields[1].name != "Extended (1/2)" {
		t.Fatalf("unexpected split field name: %q", fields[1].name)
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
}

func TestValidateSlackCfg(t *testing.T) {
	if err := validateSlackCfg(cmdCfg{}); err == nil {
		t.Fatal("expected error for empty slack config")
	}
}

func TestValidateDiscordCfg(t *testing.T) {
	if err := validateDiscordCfg(cmdCfg{}); err == nil {
		t.Fatal("expected error for empty discord config")
	}
}
