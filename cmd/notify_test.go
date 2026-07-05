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
