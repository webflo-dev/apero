package cmd

import (
	"log"
	"webflo-dev/apero-ctl/ipc"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "apply-css",
		Short: "Apply CSS",
		Run: func(cmd *cobra.Command, args []string) {
			applyCSS()
		},
	}

	rootCmd.AddCommand(cmd)
}

func applyCSS() {
	_, err := ipc.NewClient().ApplyCSS()
	if err != nil {
		log.Fatal("error:", err)
	}
}
