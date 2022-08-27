package pkg

import (
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func readFile(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file '%s'", filename)
	}

	return buf, nil
}

func writeStringToFile(filename string, text string) error {
	if err := os.WriteFile(filename, []byte(text), fs.FileMode(os.O_WRONLY)); err != nil {
		return errors.Wrapf(err, "failed to write to file '%s'", filename)
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
		return errors.Wrapf(err, "failed to parse config from '%s'", filename)
	}

	return nil
}
