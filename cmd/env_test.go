package cmd

import "testing"

func TestLoadEnv(t *testing.T) {
	t.Setenv("NAME", "Test Program")
	t.Setenv("CHANNELNAME", "NHK")
	t.Setenv("CHANNELTYPE", "GR")
	t.Setenv("STARTAT", "1728000000000")
	t.Setenv("ENDAT", "1728003600000")
	t.Setenv("ERROR_CNT", "3")

	env, err := loadEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Name != "Test Program" {
		t.Fatalf("unexpected name: %q", env.Name)
	}
	if env.ChannelName != "NHK" || env.ChannelType != "GR" {
		t.Fatalf("unexpected channel: %q/%q", env.ChannelName, env.ChannelType)
	}
	if env.StartAt != 1728000000000 || env.EndAt != 1728003600000 {
		t.Fatalf("unexpected time: start=%d end=%d", env.StartAt, env.EndAt)
	}
	if env.ErrorCnt != "3" {
		t.Fatalf("unexpected error count: %q", env.ErrorCnt)
	}
}

func TestLoadEnvInvalidInt(t *testing.T) {
	t.Setenv("STARTAT", "not-an-int")
	if _, err := loadEnv(); err == nil {
		t.Fatal("expected error when STARTAT is not an integer")
	}
}

// loadCfg reads config.yml next to the executable. During tests the binary
// lives in a temp directory with no config.yml, so the load must fail cleanly.
func TestLoadCfgMissingFile(t *testing.T) {
	if _, err := loadCfg(); err == nil {
		t.Fatal("expected error when config.yml is absent next to the test binary")
	}
}
