package pkg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

type RunOpts struct {
	DryRun     bool
	Confirm    bool
	Publish    bool
	AllowDirty bool
	Force      bool
	BaseBranch string
}

func NewRunOpts(dryRun, confirm, publish, allowDirty, force bool, baseBranch string) *RunOpts {
	return &RunOpts{
		DryRun:     dryRun,
		Confirm:    confirm,
		Publish:    publish,
		AllowDirty: allowDirty,
		BaseBranch: baseBranch,
		Force:      force,
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

	if !runOpts.Force && len(updatedTargetPaths) == 0 {
		return nil
	}

	if !runOpts.DryRun {
		for _, hook := range config.Hooks.Post {
			fmt.Printf("Running post hook '%s'\n", hook)

			cmdParts := strings.Split(hook, " ")
			args := []string{}
			if len(cmdParts) > 0 {
				args = append(args, cmdParts[1:]...)
			}

			outBytes, err := exec.CommandContext(ctx, cmdParts[0], args...).CombinedOutput() // nolint:gosec
			out := string(outBytes)
			if err != nil {
				return errors.Wrapf(err, "Failed to run post hook '%s': %s", hook, out)
			}

			if out != "" {
				fmt.Printf("Output: %s\n", out)
			}
		}
	}

	if !runOpts.DryRun && runOpts.Publish {
		if !runOpts.Confirm && !utils.AskUserYesNoQuestion("Do you want to publish the above changes?") {
			return errors.New("User aborted")
		}

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

	params := map[string]interface{}{}
	for _, paramsPath := range target.ParamsPaths {
		var curParams map[string]interface{}
		if err := utils.ReadYaml(paramsPath, &curParams); err != nil {
			return false, errors.Wrap(err, "Failed to parse params")
		}
		params = lo.Assign(params, curParams)
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
