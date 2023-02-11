package pkg

import (
	"context"
	"os"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"
	"github.com/pkg/fileutils"
	"github.com/samber/lo"

	"github.com/ilaif/goplicate/pkg/config"
	"github.com/ilaif/goplicate/pkg/git"
	"github.com/ilaif/goplicate/pkg/utils"
)

func RunTarget(ctx context.Context, target config.Target, cloner git.Cloner, dryRun, confirm bool) (bool, error) {
	workdir := utils.MustGetwd()

	sourcePath, err := ResolveSourcePath(ctx, target.Source, workdir, cloner)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to resolve source '%s'", target.Source.String())
	}

	if target.SyncInitial {
		if _, err := os.Stat(target.Path); errors.Is(err, os.ErrNotExist) {
			log.Infof("Syncing initial state of '%s' from '%s'", target.Path, sourcePath)
			if err := fileutils.CopyFile(target.Path, sourcePath); err != nil {
				return false, errors.Wrapf(err, "Failed to copy '%s' to '%s'", sourcePath, target.Path)
			}
		}
	}

	targetBlocks, err := parseBlocksFromFile(target.Path, nil)
	if err != nil {
		return false, errors.Wrap(err, "Failed to parse target blocks")
	}

	params := map[string]interface{}{}
	for _, paramsSource := range target.Params {
		paramsPath, err := ResolveSourcePath(ctx, paramsSource, workdir, cloner)
		if err != nil {
			return false, errors.Wrapf(err, "Failed to resolve source '%s'", paramsSource.String())
		}

		var curParams map[string]interface{}
		if err := utils.ReadYaml(paramsPath, &curParams); err != nil {
			return false, errors.Wrap(err, "Failed to parse params")
		}
		params = lo.Assign(params, curParams)
	}

	sourceBlocks, err := parseBlocksFromFile(sourcePath, params)
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
			log.Warnf("Target '%s': Block '%s' not found. Skipping", target.Path, targetBlock.Name)

			continue
		}

		diff := targetBlock.Compare(sourceBlock.Lines)
		if diff != "" {
			log.Infof("Target '%s': Block '%s' needs to be updated. Diff:\n%s\n", target.Path, targetBlock.Name, diff)

			targetBlock.SetLines(sourceBlock.Lines)
			anyDiff = true
		}
	}

	if !anyDiff {
		return false, nil
	}

	if dryRun {
		log.Infof("Target '%s': In dry-run mode - Not performing any changes", target.Path)

		return false, nil
	}

	question := "Do you want to apply the above changes?"
	answer, err := utils.PromptUserYesNoQuestion(question, confirm)
	if err != nil {
		return false, err
	}

	if answer {
		if err := utils.WriteStringToFile(target.Path, targetBlocks.Render()); err != nil {
			return false, err
		}

		log.Infof("Target '%s': Updated", target.Path)
	} else {
		log.Infof("Target '%s': Skipped", target.Path)
	}

	return true, nil
}
