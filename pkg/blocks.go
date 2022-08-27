package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	ParamName = "name"
	ParamPos  = "pos"

	PosStart = "start"
	PosEnd   = "end"
)

var (
	PosList = []string{PosStart, PosEnd}
)

var (
	annotationIdentifier = "goplicate"
	commentRegexp        = `(#|\/\/|\/\*|\-\-)`
	blockRegex           = regexp.MustCompile(fmt.Sprintf(`\s*%s\s*%s\((.*)\)`, commentRegexp, annotationIdentifier))
)

type Block struct {
	Name  string
	Lines []string
}

func (b *Block) Render() string {
	return strings.Join(b.Lines, "\n")
}

func (b *Block) SetLines(lines []string) {
	b.Lines = lines
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

func parseBlocksFromFile(filename string) (Blocks, error) {
	targetPathBytes, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	targetLines := splitLines(targetPathBytes)
	targetBlocks, err := parseBlocksFromLines(targetLines)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse target blocks in '%s'", filename)
	}

	return targetBlocks, err
}

func parseBlocksFromLines(lines []string) (Blocks, error) {
	blocks := Blocks{}

	var startI int
	var curBlock *Block
	for i, l := range lines {
		matches := blockRegex.FindStringSubmatch(l)
		if len(matches) != 3 {
			// if not a block comment, open an empty block
			if curBlock == nil {
				curBlock = &Block{Name: ""}
			}

			continue
		}

		params, err := parseBlockParams(matches[2])
		if err != nil {
			return nil, err
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
			curBlock.Lines = lines[startI:i]
			blocks.add(curBlock)
			startI = i
			curBlock = nil
		default:
			return nil, errors.Errorf("every block must have a position and cannot be nested"+
				"or interleaved with other blocks. params: '%s'", params)
		}
	}

	if curBlock != nil {
		if curBlock.Name == "" {
			curBlock.Lines = lines[startI:]
			blocks.add(curBlock)
		} else {
			return nil, errors.Errorf("every block must have an 'end' position")
		}
	}

	return blocks, nil
}

type blockParams struct {
	name string
	pos  string
}

func parseBlockParams(params string) (*blockParams, error) {
	bp := &blockParams{}

	for _, p := range strings.Split(params, ",") {
		splitP := strings.Split(p, "=")
		if len(splitP) != 2 {
			return nil, errors.Errorf("block parameter '%s' is not of the form 'name=value'", p)
		}
		paramName := splitP[0]
		paramValue := splitP[1]
		switch paramName {
		case "name":
			bp.name = paramValue
		case "pos":
			bp.pos = paramValue
		default:
			return nil, errors.Errorf("unknown block parameter name '%s'", p)
		}
	}

	if bp.name == "" {
		return nil, errors.Errorf("block parameter 'name' cannot be empty")
	}

	if !lo.Contains(PosList, bp.pos) {
		return nil, errors.Errorf("block parameter 'pos' must be one of %s", PosList)
	}

	return bp, nil
}
