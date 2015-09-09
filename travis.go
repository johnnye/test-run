package main

import (
	"gopkg.in/yaml.v2"
	"errors"
	"strings"
)

type Travis struct {
	Filename string
	Script   string
}

func (t *Travis) filename() string {
	if t.Filename == "" {
		t.Filename = ".travis.yml"
	}
	return t.Filename
}

func (t *Travis) runTests() error {
	err := errors.New("")

	if !doesFileExist(t.filename()) {
		s := []string{t.filename(), "file does not exist"};
		err = errors.New(strings.Join(s, " "))
		return err
	}
	raw := readFile(t.filename())

	err = t.getCommandsFromYAML([]byte(raw))

	err = executeCommands(t.Script)

	return err
}

func (t *Travis) getCommandsFromYAML(raw []byte) error {
	err := yaml.Unmarshal(raw, &t)
	return err
}



