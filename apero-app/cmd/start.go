package cmd

import (
	ipc "apero-ipc"
	"apero/app"
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
