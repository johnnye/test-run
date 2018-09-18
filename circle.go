package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"strings"
)

type Circle struct {
	Filename string
	Test     struct{ Override []string }
}

func (c *Circle) filename() string {
	if c.Filename == "" {
		c.Filename = "circle.yml"
	}
	return c.Filename
}

func (c *Circle) runTests() error {
	err := errors.New("")

	if !doesFileExist(c.filename()) {
		s := []string{c.filename(), "file does not exist"}
		err = errors.New(strings.Join(s, " "))
		return err
	}

	raw := readFile(c.filename())

	err = c.getCommandsFromYAML([]byte(raw))

	for _, cmd := range c.Test.Override {
		err = executeCommands(cmd)
	}

	return err
}

func (c *Circle) getCommandsFromYAML(raw []byte) error {
	err := yaml.Unmarshal(raw, &c)
	return err
}
