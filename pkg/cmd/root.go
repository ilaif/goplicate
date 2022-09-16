package cmd

import (
	"context"
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

var rootCmd = &cobra.Command{
	Use:              "goplicate",
	Short:            "Sync code or configuration snippets from a source repository to multiple target projects",
	SilenceUsage:     true,
	PersistentPreRun: toggleDebug,
}

func Execute() {
	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose logging")
}

func toggleDebug(cmd *cobra.Command, args []string) {
	log.SetLevel(log.InfoLevel)
	log.DecreasePadding() // remove the default padding

	if debug {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
	}
}
