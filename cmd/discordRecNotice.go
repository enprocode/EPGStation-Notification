package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var discordRecStartCmd = &cobra.Command{
	Use:   "discordRecStart",
	Short: "Recording start notification command",
	Long:  `This command notifies the start of recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DiscordSend(" :arrow_forward: ", 3447003, false); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var discordRecEndCmd = &cobra.Command{
	Use:   "discordRecEnd",
	Short: "Recording end notification command",
	Long:  `This command notifies the end of recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DiscordSend(" :white_check_mark: ", 3066993, false); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var discordRecErrorCmd = &cobra.Command{
	Use:   "discordRecError",
	Short: "Recording error notification command",
	Long:  `This command notifies you of recording errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DiscordSend(" :x: ", 15158332, true); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(discordRecStartCmd)
	rootCmd.AddCommand(discordRecEndCmd)
	rootCmd.AddCommand(discordRecErrorCmd)
}
