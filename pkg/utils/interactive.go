package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/pkg/errors"
)

// PromptUserYesNoQuestion ask a question and wait for user input.
func PromptUserYesNoQuestion(question string, confirm bool) (bool, error) {
	if confirm {
		return true, nil
	}

	var continueToAuth bool

	if err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &continueToAuth); err != nil {
		if err == terminal.InterruptErr {
			return false, errors.Wrap(err, "user interrupt")
		}

		return false, errors.Wrap(err, "prompt error")
	}

	return continueToAuth, nil
}

// OpenTextEditor opens the default text editor and capturing its input.
func OpenTextEditor(ctx context.Context, initMsg string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	f, err := ioutil.TempFile("", "go-editor")
	if err != nil {
		return "", errors.Wrap(err, "Failed to create temp file")
	}
	defer os.Remove(f.Name())

	if initMsg != "" {
		if _, err := f.WriteString(initMsg); err != nil {
			return "", errors.Wrap(err, "Failed to write init message to temp file")
		}
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("%s %s", editor, f.Name())) // nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Wrap(err, "Failed to run text editor")
	}

	bytes, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", errors.Wrap(err, "Failed to read temp file")
	}

	return string(bytes), nil
}
