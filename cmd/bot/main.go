package bot

import "github.com/spf13/cobra"

func NewBotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "linebot [service]",
		Short: "start comment's service",
	}

	cmd.AddCommand(newBotCommand())
	cmd.AddCommand(newMigrationCommand())

	return cmd
}
