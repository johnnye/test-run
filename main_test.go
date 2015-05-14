package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

var goodData = `
test:
  override:
    - ./vendor/bin/parallel-lint --exclude vendor .
    - ./vendor/bin/phpunit -c phpunit.xml.dist
    - ./vendor/bin/phpcs --standard=vendor/crowdcube/codesniffer-standard/Crowdcube -p --report=full --report-checkstyle=build/logs/checkstyle.xml --runtime-set ignore_warnings_on_exit true src/ lib/ tests/
`

var simpleData = `
test:
  override:
    - ./
    - ./v
    - ./ve
`

func TestUnderstandsYAML(t *testing.T) {
	Convey("test reading and running tests", t, func() {

		ioutil.WriteFile("/tmp/circle.yml", []byte(goodData), 0644)

		Convey("Should find testing -> override in YMAL", func() {
			t := Circle{}
			t.Test.Command = []string{"./", "./v", "./ve"}

			r, err := getCommandsFromYAML([]byte(simpleData))

			So(r, ShouldResemble, t)
			So(err, ShouldBeEmpty)
		})

		Convey("Should log an error with bad YAML", func() {

			t, err := getCommandsFromYAML([]byte("sdkfjls"))

			So(err, ShouldNotBeNil)
			So(t, ShouldResemble, Circle{})
		})

		Convey("Circle.yml should not exist", func() {
			So(doesACircleFileExist("/tmp/baz.yml"), ShouldEqual, false)
		})

		Convey("circle.yml should exist", func() {
			So(doesACircleFileExist("/tmp/circle.yml"), ShouldEqual, true)
		})

		Convey("Finding a circlefile should be good!", func() {
			So(readCircleFile("/tmp/circle.yml"), ShouldEqual, goodData)
		})

		Convey("clean vendor bin from string", func() {
			So(cleanVendorBin("./vendor/bin/baz"), ShouldEqual, "./baz")
		})

		Convey("no need to clean vendor bin", func() {
			So(cleanVendorBin("./bin/baz"), ShouldEqual, "./bin/baz")
		})

		Convey("run and echo a pwd command", func() {
			c := "pwd"
			err := executeCommands(c)
			So(err, ShouldBeNil)
		})

		Convey("run a command with invalid arguments", func() {
			So(executeCommands("pwd -e"), ShouldNotBeNil)
		})

		Convey("run all of the above", func() {
			So(runCircleTests("/tmp/curcle.yml"), ShouldBeNil)
		})

		Reset(func() {
			os.Remove("/tmp/circle.yml")
		})
	})
}
