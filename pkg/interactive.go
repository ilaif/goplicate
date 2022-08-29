package pkg

import (
	"github.com/AlecAivazis/survey/v2"
)

// askUserYesNoQuestion ask a question and wait for user input.
func askUserYesNoQuestion(question string) bool {
	var continueToAuth bool

	if err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &continueToAuth); err != nil {
		return false
	}

	return continueToAuth
}
