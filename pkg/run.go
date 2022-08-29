package pkg

import (
	"fmt"

	"github.com/pkg/errors"
)

type RunOpts struct {
	DryRun  bool
	Confirm bool
}

func Run(config *ProjectConfig, runOpts RunOpts) error {
	for _, target := range config.Targets {
		if err := runTarget(target, runOpts); err != nil {
			return errors.Wrapf(err, "Target '%s'", target.Path)
		}
	}

	return nil
}

func runTarget(target Target, runOpts RunOpts) error {
	targetBlocks, err := parseBlocksFromFile(target.Path, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to parse target blocks")
	}

	targetSource, err := parseTargetSource(target.Source)
	if err != nil {
		return err
	}

	var params map[string]interface{}
	if err := readYaml(target.ParamsPath, &params); err != nil {
		return errors.Wrap(err, "Failed to parse params")
	}

	sourceBlocks, err := parseBlocksFromFile(targetSource.Path, params)
	if err != nil {
		return errors.Wrap(err, "Failed to parse source blocks")
	}

	anyDiff := false

	for _, targetBlock := range targetBlocks {
		if targetBlock.Name == "" {
			continue
		}

		sourceBlock := sourceBlocks.Get(targetBlock.Name)
		if sourceBlock == nil {
			fmt.Printf("WARNING: Target '%s': Block '%s' not found. Skipping\n", target.Path, targetBlock.Name)

			continue
		}

		diff := targetBlock.Compare(sourceBlock.Lines)
		if diff != "" {
			fmt.Printf("Target '%s': Block '%s' needs to be updated. Diff:\n\n", target.Path, targetBlock.Name)
			fmt.Printf("%s\n\n", diff)

			targetBlock.SetLines(sourceBlock.Lines)
			anyDiff = true
		}
	}

	if !anyDiff {
		return nil
	}

	if runOpts.DryRun {
		fmt.Printf("Target '%s': In dry-run mode - Not performing any changes\n", target.Path)

		return nil
	}

	if !runOpts.Confirm && !askUserYesNoQuestion("Do you want to apply the above changes?") {
		return nil
	}

	if err := writeStringToFile(target.Path, targetBlocks.Render()); err != nil {
		return err
	}

	fmt.Printf("Target '%s': Updated\n", target.Path)

	return nil
}
