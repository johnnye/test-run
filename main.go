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
)

type Tests struct {
	Test struct{ Override []string }
}

type CmdArg struct {
	Cmd  string
	args []string
}

func main() {
	filename := "circle.yml"

	if !doesACircleFileExist(filename) {
		//TODO: something here
	}

	raw := readCircleFile(filename)

	commands := getCommandsFromYAML([]byte(raw))

	for _, cmd := range commands.Test.Override {

		executeCommands(cmd)
	}
}

func getCommandsFromYAML(raw []byte) (tests Tests) {
	t := Tests{}

	err := yaml.Unmarshal(raw, &t)

	if err != nil {
		log.Fatal("error cannot parse YAML")
	}
	return t
}

func readCircleFile(filename string) (contents string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {

	}

	return string(data)
}

func doesACircleFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", filename)
		return false
	}
	return true
}

func cleanVendorBin(unclean string) (clean string) {
	return strings.Replace(unclean, "./vendor/bin/", "./", 1)
}

func executeCommands(incoming string) {

	//
	// Taken from http://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	//

	cmd := exec.Command("bash", "-c", incoming)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
}
