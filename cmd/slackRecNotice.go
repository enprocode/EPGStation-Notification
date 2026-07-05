package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var slackRecStartCmd = &cobra.Command{
	Use:   "slackRecStart",
	Short: "Recording start notification command",
	Long:  `This command notifies the start of recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Slack(" :arrow_forward: ", "#439FE0", false); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var slackRecEndCmd = &cobra.Command{
	Use:   "slackRecEnd",
	Short: "Recording end notification command",
	Long:  `This command notifies the end of recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Slack(" :white_check_mark: ", "good", false); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var slackRecErrorCmd = &cobra.Command{
	Use:   "slackRecError",
	Short: "Recording error notification command",
	Long:  `This command notifies you of recording errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Slack(" :x: ", "danger", true); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(slackRecStartCmd)
	rootCmd.AddCommand(slackRecEndCmd)
	rootCmd.AddCommand(slackRecErrorCmd)
}
