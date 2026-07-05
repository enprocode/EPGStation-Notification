package cmd

import (
	"epgstation_notification/cmd"
	"testing"
)

func TestDiscordSend(t *testing.T) {
	type args struct {
		Icon          string
		Col           int
		WithErrorInfo bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.DiscordSend(tt.args.Icon, tt.args.Col, tt.args.WithErrorInfo); (err != nil) != tt.wantErr {
				t.Errorf("DiscordSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
