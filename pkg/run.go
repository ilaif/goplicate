package pkg

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

type RunOpts struct {
	DryRun     bool
	Confirm    bool
	Publish    bool
	AllowDirty bool
	BaseBranch string
}

func NewRunOpts(dryRun, confirm, publish, allowDirty bool, baseBranch string) *RunOpts {
	return &RunOpts{
		DryRun:     dryRun,
		Confirm:    confirm,
		Publish:    publish,
		AllowDirty: allowDirty,
		BaseBranch: baseBranch,
	}
}

func Run(ctx context.Context, config *ProjectConfig, runOpts *RunOpts) error {
	publisher := git.NewPublisher(runOpts.BaseBranch)

	if !runOpts.DryRun && runOpts.Publish {
		if err := publisher.Init(ctx, "."); err != nil {
			return errors.Wrap(err, "Failed to initialize git")
		}

		if !runOpts.AllowDirty && !publisher.IsClean() {
			return errors.New("Git worktree is not clean. Please commit or stash changes before running again")
		}
	}

	updatedTargetPaths := []string{}
	for _, target := range config.Targets {
		if updated, err := runTarget(target, runOpts); err != nil {
			return errors.Wrapf(err, "Target '%s'", target.Path)
		} else if updated {
			updatedTargetPaths = append(updatedTargetPaths, target.Path)
		}
	}

	if !runOpts.DryRun && runOpts.Publish && len(updatedTargetPaths) > 0 {
		if err := publisher.Publish(ctx, updatedTargetPaths); err != nil {
			return errors.Wrap(err, "Failed to publish changes")
		}
	}

	return nil
}

func runTarget(target Target, runOpts *RunOpts) (bool, error) {
	targetBlocks, err := parseBlocksFromFile(target.Path, nil)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse target blocks")
	}

	targetSource, err := parseTargetSource(target.Source)
	if err != nil {
		return false, err
	}

	var params map[string]interface{}
	if err := utils.ReadYaml(target.ParamsPath, &params); err != nil {
		return false, errors.Wrap(err, "Failed to parse params")
	}

	sourceBlocks, err := parseBlocksFromFile(targetSource.Path, params)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse source blocks")
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
		return false, nil
	}

	if runOpts.DryRun {
		fmt.Printf("Target '%s': In dry-run mode - Not performing any changes\n", target.Path)

		return false, nil
	}

	if !runOpts.Confirm && !utils.AskUserYesNoQuestion("Do you want to apply the above changes?") {
		return false, errors.New("User aborted")
	}

	if err := utils.WriteStringToFile(target.Path, targetBlocks.Render()); err != nil {
		return false, err
	}

	fmt.Printf("Target '%s': Updated\n", target.Path)

	return true, nil
}
