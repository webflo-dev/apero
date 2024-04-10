package cmd

import (
	"apero/app"
	"apero/ipc"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			ipc.StartIPC()
			exitCode := app.Start()

			os.Exit(exitCode)
		},
	}

	rootCmd.AddCommand(cmd)
}
