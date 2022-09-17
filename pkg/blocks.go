package pkg

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/ilaif/goplicate/pkg/utils"
)

const (
	ParamName = "name"
	ParamPos  = "pos"

	PosStart = "start"
	PosEnd   = "end"
)

var (
	PosList = []string{PosStart, PosEnd}

	// regex decomposition:
	// 1. empty spaces: \s*
	// 2. identifying comments: (#|\/\/|\/\*|\-\-|<\-\-)
	// 3. goplicate block format: goplicate_start|end(...params...) or goplicate-start:<name>
	blockRegex = regexp.MustCompile(`\s*(#|\/\/|\/\*|\-\-|<\-\-)\s*goplicate([_\-](start|end))?(\((.*)\)|:(.*))`)
)

type Block struct {
	Name  string
	Lines []string
}

func (b *Block) Render() string {
	return strings.Join(b.Lines, "\n")
}

func (b *Block) Compare(lines []string) string {
	return linesDiff(b.Lines, b.padLines(lines))
}

func (b *Block) SetLines(lines []string) {
	b.Lines = b.padLines(lines)
}

// padLines add a base indentation to match the one in this block (according to the first line)
func (b *Block) padLines(lines []string) []string {
	ourIndent := 0
	if len(b.Lines) > 0 {
		ourIndent = utils.CountLeadingSpaces(b.Lines[0])
	}
	theirIndent := 0
	if len(lines) > 0 {
		theirIndent = utils.CountLeadingSpaces(lines[0])
	}
	indentAddition := ourIndent - theirIndent

	paddedLines := make([]string, len(lines))
	for i, l := range lines {
		if indentAddition > 0 {
			paddedLines[i] = strings.Repeat(" ", indentAddition) + l
		} else {
			paddedLines[i] = l[-indentAddition:]
		}
	}

	return paddedLines
}

type Blocks []*Block

func (b *Blocks) add(block *Block) {
	*b = append(*b, block)
}

func (b *Blocks) Get(name string) *Block {
	if name == "" {
		return nil
	}

	for _, block := range *b {
		if block.Name == name {
			return block
		}
	}

	return nil
}

func (b *Blocks) Render() string {
	return strings.Join(lo.Map(*b, func(block *Block, _ int) string {
		return block.Render()
	}), "\n")
}

func parseBlocksFromFile(filename string, params map[string]interface{}) (Blocks, error) {
	fileBytes, err := utils.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var s string
	if params != nil {
		t, err := template.New("parse-blocks-tpl").Parse(string(fileBytes))
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to parse template for file '%s'", filename)
		}

		var tpl bytes.Buffer
		if err := t.Option("missingkey=error").Execute(&tpl, params); err != nil {
			return nil, errors.Wrapf(err, "Failed to execute template for file '%s'", filename)
		}

		s = tpl.String()
	} else {
		s = string(fileBytes)
	}

	lines := strings.Split(s, "\n")

	blocks, err := parseBlocksFromLines(lines)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse blocks in '%s'", filename)
	}

	return blocks, err
}

func parseBlocksFromLines(lines []string) (Blocks, error) {
	blocks := Blocks{}

	var startI int
	var curBlock *Block
	for i, l := range lines {
		params, err := parseBlockComment(l)
		if err != nil {
			return nil, err
		} else if params == nil {
			// if not a block comment, open an empty block
			if curBlock == nil {
				curBlock = &Block{Name: ""}
			}

			continue
		}

		if params.pos == PosStart && curBlock != nil && curBlock.Name == "" {
			// an empty block should be closed before processing the comment
			curBlock.Lines = lines[startI:i]
			blocks.add(curBlock)
			startI = i
			curBlock = nil
		}

		switch {
		case params.pos == PosStart && curBlock == nil:
			// if we see a block start with a nil curBlock, we'll initialize a new one
			curBlock = &Block{Name: params.name}
		case params.pos == PosEnd && curBlock != nil:
			// if we see a block end with a currently active curBlock, we'll close it
			curBlock.Lines = lines[startI : i+1]
			blocks.add(curBlock)
			startI = i + 1
			curBlock = nil
		default:
			return nil, errors.Errorf("Every block must have a position and cannot be nested "+
				"or interleaved with other blocks. Params: '%+v'", params)
		}
	}

	if curBlock != nil {
		if curBlock.Name == "" {
			curBlock.Lines = lines[startI:]
			blocks.add(curBlock)
		} else {
			return nil, errors.Errorf("Every block must have an 'end' position")
		}
	}

	return blocks, nil
}

type blockParams struct {
	name string
	pos  string
}

func parseBlockComment(l string) (*blockParams, error) {
	matches := blockRegex.FindStringSubmatch(l)
	if len(matches) != 7 {
		// if not a block comment, return nil
		return nil, nil
	}

	startEndBlock := matches[3]
	paramsStr := matches[5]
	if matches[6] != "" {
		// Assume the format is "goplicate-start:<name>""
		paramsStr = fmt.Sprintf("name=%s", matches[6])
	}

	return parseBlockParams(startEndBlock, paramsStr)
}

func parseBlockParams(startEndBlock string, params string) (*blockParams, error) {
	bp := &blockParams{}

	if startEndBlock != "" {
		bp.pos = startEndBlock
	}

	for _, p := range strings.Split(params, ",") {
		splitP := strings.Split(p, "=")
		if len(splitP) != 2 {
			return nil, errors.Errorf("Block parameter '%s' is not of the form 'name=value'", p)
		}
		paramName := splitP[0]
		paramValue := splitP[1]
		switch paramName {
		case "name":
			bp.name = paramValue
		case "pos":
			bp.pos = paramValue
		default:
			return nil, errors.Errorf("Unknown block parameter name '%s'", p)
		}
	}

	if bp.name == "" {
		return nil, errors.Errorf("Block parameter 'name' cannot be empty")
	}

	if !lo.Contains(PosList, bp.pos) {
		return nil, errors.Errorf("Block parameter 'pos' must be one of %s", PosList)
	}

	return bp, nil
}
