package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"errors"
)

//
// Tests is used to unmarshal the yaml
// Overrides is an array of the tests to run
//
type Tests struct {
	Test struct{ Override []string }
}

func main() {
	filename := "circle.yml"

	err := runCircleTests(filename)

	if err != nil {
		log.Println(err)
	}
}

func runCircleTests(file string) (error) {
	err := errors.New("")

	if !doesACircleFileExist(file) {
		err = errors.New("file does not exist")
	}

	raw := readCircleFile(file)

	commands, err := getCommandsFromYAML([]byte(raw))

	for _, cmd := range commands.Test.Override {
		err = executeCommands(cmd)
	}
	return err
}

func getCommandsFromYAML(raw []byte) (Tests, error) {
	t := Tests{}
	err := yaml.Unmarshal(raw, &t)
	return t, err
}

func readCircleFile(filename string) (contents string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {

	}

	return string(data)
}

func doesACircleFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func cleanVendorBin(unclean string) (clean string) {
	return strings.Replace(unclean, "./vendor/bin/", "./", 1)
}

func executeCommands(incoming string) (error) {

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
