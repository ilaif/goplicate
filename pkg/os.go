package pkg

import (
	"os"

	"github.com/pkg/errors"
)

func MustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return wd
}

func Chdir(dir string) error {
	if err := os.Chdir(dir); err != nil {
		return errors.Wrap(err, "Failed to change directory")
	}

	return nil
}
