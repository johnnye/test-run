package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"os"
	"gopkg.in/yaml.v2"
	"log"
)

type Provider interface {
	filename() string
	runTests() error
	getCommandsFromYAML([]byte) error
}

//
// Circle is used to unmarshal the yaml
// Overrides is an array of the tests to run
//
type Circle struct {
	Filename string
	Test     struct{ Command []string }
}

type Travis struct {
	Filename string
	Script   string
}

func main() {

	providers := []Provider{&Circle{}, &Travis{}}
	err := errors.New("")
	for _, provider := range providers {
		err = provider.runTests()
		fmt.Println("===================")
	}
	log.Println(err)
}

func cleanVendorBin(unclean string) (clean string) {
	return strings.Replace(unclean, "./vendor/bin/", "./", 1)
}

func readFile(filename string) (contents string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {

	}

	return string(data)
}

func executeCommands(incoming string) error {

	//
	// Taken from http://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	//

	cmd := exec.Command("bash", "-c", incoming)
	cmdReader, err := cmd.StdoutPipe()

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	err = cmd.Wait()
	return err
}

func doesFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
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
		err = errors.New("file does not exist")
	}

	raw := readFile(c.filename())

	err = c.getCommandsFromYAML([]byte(raw))

	for _, cmd := range c.Test.Command {
		err = executeCommands(cmd)
	}

	return err
}

func (c *Circle) getCommandsFromYAML(raw []byte) error {
	err := yaml.Unmarshal(raw, &c)
	return err
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
		err = errors.New("file does not exist")
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



