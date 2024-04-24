package cmd

import (
	"log"

	"webflo-dev/apero-ctl/ipc"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "quit",
		Short: "Quit apero :'(",
		Run: func(cmd *cobra.Command, args []string) {
			sendQuit()
		},
	}

	rootCmd.AddCommand(cmd)
}

func sendQuit() {
	_, err := ipc.NewClient().SendQuit()
	if err != nil {
		log.Fatal("error:", err)
	}
}
