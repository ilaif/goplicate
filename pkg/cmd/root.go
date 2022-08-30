package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "goplicate",
	Short:        "Sync code or configuration snippets from a source repository to multiple target projects",
	SilenceUsage: true,
}

func Execute() {
	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
}
