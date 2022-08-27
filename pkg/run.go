package pkg

import (
	"github.com/pkg/errors"
)

func Run(config *RepositoryConfig) error {
	for _, target := range config.Targets {
		targetBlocks, err := parseBlocksFromFile(target.Path, nil)
		if err != nil {
			return errors.Wrapf(err, "failed to parse target blocks in '%s'", target.Path)
		}

		targetSource, err := parseTargetSource(target.Source)
		if err != nil {
			return err
		}

		var params map[string]interface{}
		if err := readYaml(target.ParamsPath, &params); err != nil {
			return errors.Wrap(err, "failed to parse params")
		}

		sourceBlocks, err := parseBlocksFromFile(targetSource.Path, params)
		if err != nil {
			return errors.Wrapf(err, "failed to parse source blocks in '%s'", target.Path)
		}

		for _, targetBlock := range targetBlocks {
			sourceBlock := sourceBlocks.Get(targetBlock.Name)
			if sourceBlock == nil {
				continue
			}

			targetBlock.SetLines(sourceBlock.Lines)
		}

		if err := writeStringToFile(target.Path, targetBlocks.Render()); err != nil {
			return err
		}
	}

	return nil
}
