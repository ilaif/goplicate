package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ilaif/goplicate/pkg/utils"
)

func TestParseBlocksFromFile(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		file           string
		error          bool
		expectedBlocks Blocks
	}{
		{
			file:  "testdata/blocks/valid.yaml",
			error: false,
			expectedBlocks: Blocks{
				{
					Name: "",
					Lines: []string{
						"repos:",
					},
				},
				{
					Name: "common",
					Lines: []string{
						"  # goplicate(name=common,pos=start)",
						"  - name: common",
						"    hooks:",
						"      - id: my-common-pre-commit-hook",
						"  # goplicate(name=common,pos=end)",
					},
				},
				{
					Name: "",
					Lines: []string{
						"  - name: local",
						"    hooks:",
						"      - id: my-project-1-pre-commit-hook",
					},
				},
				{
					Name: "external",
					Lines: []string{
						"  # goplicate_start(name=external)",
						"  - name: external",
						"    hooks:",
						"      - id: my-external-pre-commit-hook",
						"  # goplicate_end(name=external)",
					},
				},
				{
					Name: "new",
					Lines: []string{
						"  # goplicate-start:new",
						"  - name: new",
						"    hooks:",
						"      - id: my-new-pre-commit-hook",
						"  # goplicate-end:new",
					},
				},
				{
					Name: "",
					Lines: []string{
						"",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			blocks, err := parseBlocksFromFile(test.file, nil)
			a.NoError(err)

			a.Equal(test.expectedBlocks, blocks)
		})
	}
}

func TestBlocksPadding(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		targetBlock Block
		sourceBlock Block
	}{
		{
			targetBlock: Block{
				Name: "common",
				Lines: []string{
					"  # goplicate-start:common",
					"  value",
					"  # goplicate-end:common",
				},
			},
			sourceBlock: Block{
				Name: "common",
				Lines: []string{
					"    # goplicate-start:common",
					"    value",
					"    # goplicate-end:common",
				},
			},
		},
		{
			targetBlock: Block{
				Name: "common",
				Lines: []string{
					"    # goplicate-start:common",
					"    value",
					"    # goplicate-end:common",
				},
			},
			sourceBlock: Block{
				Name: "common",
				Lines: []string{
					"  # goplicate-start:common",
					"  value",
					"  # goplicate-end:common",
				},
			},
		},
		{
			targetBlock: Block{
				Name: "common",
				Lines: []string{
					"  # goplicate-start:common",
					"  value",
					"  # goplicate-end:common",
				},
			},
			sourceBlock: Block{
				Name: "common",
				Lines: []string{
					"  # goplicate-start:common",
					"  value",
					"  # goplicate-end:common",
				},
			},
		},
	}

	for _, test := range tests {
		lines := test.targetBlock.padLines(test.sourceBlock.Lines)
		expectedLinePadding := utils.CountLeadingSpaces(test.targetBlock.Lines[0])
		for _, line := range lines {
			actualLinePadding := utils.CountLeadingSpaces(line)
			a.Equal(expectedLinePadding, actualLinePadding)
		}
	}
}
