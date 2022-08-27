package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ilaif/goplicate/pkg"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run goplicate on the target repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := pkg.LoadRepositoryConfig()
		if err != nil {
			return err
		}

		if err := pkg.Run(c); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
