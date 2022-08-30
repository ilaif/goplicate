package utils

import (
	"github.com/AlecAivazis/survey/v2"
)

// AskUserYesNoQuestion ask a question and wait for user input.
func AskUserYesNoQuestion(question string) bool {
	var continueToAuth bool

	if err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &continueToAuth); err != nil {
		return false
	}

	return continueToAuth
}
