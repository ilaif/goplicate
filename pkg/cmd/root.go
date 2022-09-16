package cmd

import (
	"context"
	"os"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(version string) *cobra.Command {
	var (
		debug bool
	)

	var rootCmd = &cobra.Command{
		Use:          "goplicate",
		Short:        "Sync project configuration snippets from a source repository to multiple target projects",
		SilenceUsage: true,
		Version:      version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.InfoLevel)
			log.DecreasePadding() // remove the default padding

			if debug {
				log.Info("Debug logs enabled")
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "verbose logging")

	rootCmd.AddCommand(
		newRunCmd(),
		newSyncCmd(),
	)

	return rootCmd
}
