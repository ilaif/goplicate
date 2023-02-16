package utils

import (
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
