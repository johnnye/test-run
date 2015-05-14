package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
	"os"
)

type Provider interface {
	filename() string
	runTests() error
	getCommandsFromYAML([]byte) error
	doesFileExist() bool
}

type Circle struct {
	Test struct{ Command []string }
}

func main() {

	providers := []Provider{Circle{}}
	err := errors.New("")
	for _, provider := range providers {
		err = provider.runTests()
		fmt.Println(provider.filename())
		if err != nil {
			log.Println(provider.filename(), " had this error ", err)
		}
	}
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

func (c Circle) filename() string {
	return "circle.yml"
}

func (c Circle) runTests() error {
	err := errors.New("")

	if !c.doesFileExist() {
		err = errors.New("file does not exist")
	}

	raw := readFile(c.filename())

	err = c.getCommandsFromYAML([]byte(raw))

	for _, cmd := range c.Test.Command {
		err = executeCommands(cmd)
	}

	return err
}

func (c Circle) getCommandsFromYAML(raw []byte) error {
	err := yaml.Unmarshal(raw, &c)
	return err
}

func (c Circle) doesFileExist() bool {
	if _, err := os.Stat(c.filename()); os.IsNotExist(err) {
		return false
	}
	return true
}
