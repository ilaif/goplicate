package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBlocks(t *testing.T) {
	assert := assert.New(t)

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
			assert.NoError(err)

			assert.Equal(test.expectedBlocks, blocks)
		})
	}
}
