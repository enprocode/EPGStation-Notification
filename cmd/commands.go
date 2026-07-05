package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Discord embed colors (decimal RGB) used per notification type.
const (
	colorRecStart = 3447003  // blue
	colorRecEnd   = 3066993  // green
	colorRecError = 15158332 // red
)

// recordingCommand describes a single recording-notification subcommand.
type recordingCommand struct {
	use    string
	short  string
	long   string
	notify func() error
}

// recordingCommands lists every notification subcommand. Adding a new one is a
// single entry here.
var recordingCommands = []recordingCommand{
	{"slackRecStart", "Recording start notification command", "This command notifies the start of recording.",
		func() error { return Slack(" :arrow_forward: ", "#439FE0", false) }},
	{"slackRecEnd", "Recording end notification command", "This command notifies the end of recording.",
		func() error { return Slack(" :white_check_mark: ", "good", false) }},
	{"slackRecError", "Recording error notification command", "This command notifies you of recording errors.",
		func() error { return Slack(" :x: ", "danger", true) }},
	{"discordRecStart", "Recording start notification command", "This command notifies the start of recording.",
		func() error { return DiscordSend(" :arrow_forward: ", colorRecStart, false) }},
	{"discordRecEnd", "Recording end notification command", "This command notifies the end of recording.",
		func() error { return DiscordSend(" :white_check_mark: ", colorRecEnd, false) }},
	{"discordRecError", "Recording error notification command", "This command notifies you of recording errors.",
		func() error { return DiscordSend(" :x: ", colorRecError, true) }},
}

func init() {
	for _, c := range recordingCommands {
		notify := c.notify
		rootCmd.AddCommand(&cobra.Command{
			Use:   c.use,
			Short: c.short,
			Long:  c.long,
			Run: func(cmd *cobra.Command, args []string) {
				if err := notify(); err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
			},
		})
	}
}
