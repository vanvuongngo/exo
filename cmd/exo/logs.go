package main

import (
	"github.com/deref/exo/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tails logs in an interactive terminal ui",
	Long:  `Tails logs in an interactive terminal ui`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDeamon()

		app := tui.NewApp(newClient())
		return app.Run(ctx)
	},
}
