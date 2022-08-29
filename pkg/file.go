package pkg

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func readFile(filename string) ([]byte, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get absolute path for file '%s'", filename)
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file '%s'", filename)
	}

	return buf, nil
}

func writeStringToFile(filename string, text string) error {
	if err := os.WriteFile(filename, []byte(text), fs.FileMode(os.O_WRONLY)); err != nil {
		return errors.Wrapf(err, "Failed to write to file '%s'", filename)
	}

	return nil
}

func readYaml(filename string, config interface{}) error {
	buf, err := readFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse config from '%s'", filename)
	}

	return nil
}
