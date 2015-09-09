package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"os"
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



func main() {

	providers := []Provider{&Circle{}, &Travis{}}
	err := errors.New("")
	for _, provider := range providers {
		err = provider.runTests()
		if err != nil{
			fmt.Println(err)
		}
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