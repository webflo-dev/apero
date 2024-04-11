package cmd

import (
	"fmt"
	"log"

	"webflo-dev/apero/ipc"

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
	reply, err := ipc.SendQuit()
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("AperoCtl: ", reply)
}
