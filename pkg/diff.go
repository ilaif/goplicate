package pkg

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// linesDiff get the string diff between two line slices
func linesDiff(lines1, lines2 []string) string {
	dmp := diffmatchpatch.New()

	str1dmp, str2dmp, dmpStrings := dmp.DiffLinesToChars(padLines(lines1), padLines(lines2))
	diffs := dmp.DiffMain(str1dmp, str2dmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)

	var newDiffs []diffmatchpatch.Diff

	hasDiff := false

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffDelete, diffmatchpatch.DiffInsert:
			hasDiff = true
		}

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			diff.Text = fmt.Sprintf("-%s", diff.Text[1:])
		case diffmatchpatch.DiffInsert:
			diff.Text = fmt.Sprintf("+%s", diff.Text[1:])
		}

		newDiffs = append(newDiffs, diff)
	}

	if !hasDiff {
		return ""
	}

	diff := dmp.DiffPrettyText(newDiffs)
	diffLines := strings.Split(diff, "\n")

	diffLineNos := []int{}
	for i, line := range diffLines {
		if strings.HasPrefix(line, "\x1b") {
			diffLineNos = append(diffLineNos, i)
		}
	}

	scopedLineNos := []int{}
	for _, lineNo := range diffLineNos {
		start := lo.Max([]int{lineNo - 3, 0})
		end := lo.Min([]int{lineNo + 3, len(diffLines)})
		for i := start; i < end; i++ {
			scopedLineNos = append(scopedLineNos, i)
		}
	}

	scopedLineNos = lo.Uniq(scopedLineNos)

	scopedDiff := []string{}

	if len(scopedLineNos) > 0 && scopedLineNos[0] > 0 {
		scopedDiff = append(scopedDiff, "...")
	}

	for i, lineNo := range scopedLineNos {
		if i > 0 && lineNo > scopedLineNos[i-1]+1 {
			scopedDiff = append(scopedDiff, "...")
		}

		scopedDiff = append(scopedDiff, diffLines[lineNo])
	}

	if len(scopedLineNos) > 0 && scopedLineNos[len(scopedLineNos)-1] < len(diffLines)-1 {
		scopedDiff = append(scopedDiff, "...")
	}

	return strings.Join(scopedDiff, "\n")
}

func padLines(lines []string) string {
	lines = lo.Map(lines, func(line string, _ int) string { return " " + line })

	return strings.Join(lines, "\n")
}
