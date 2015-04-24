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
			t := Tests{}
			t.Test.Override = []string{"./", "./v", "./ve"}

			So(getCommandsFromYAML([]byte(simpleData)), ShouldResemble, t)
		})

		//    SkipConvey("for the moment", func() {
		//        Convey("Should log an error with bad YAML", t, func() {
		//            So(getCommandsFromYAML([]byte("sdkfjls")), ShouldPanicWith, errors.New("error cannot parse YAML"))
		//        })
		//    })

		Convey("Circle.yml should not exist", func() {
			So(doesACircleFileExist("/tmp/baz.yml"), ShouldEqual, false)
		})

		Convey("circle.yml should exist", func() {
			So(doesACircleFileExist("/tmp/circle.yml"), ShouldEqual, true)
		})

		Convey("Finding a circlefile should be good!", func() {
			So(readCircleFile("/tmp/circle.yml"), ShouldEqual, goodData)

		})

		//		Convey("command splitter shoud split command", func() {
		//			command := "./vendor/bin/phpcs --standard=vendor/crowdcube/codesniffer-standard/Crowdcube -p --report=full --report-checkstyle=build/logs/checkstyle.xml --runtime-set ignore_warnings_on_exit true src/ lib/ tests/"
		//			cmd := "./vendor/bin/phpcs"
		//			results := []string{
		//				"--standard=vendor/crowdcube/codesniffer-standard/Crowdcube",
		//				"-p",
		//				"--report=full",
		//				"--report-checkstyle=build/logs/checkstyle.xml",
		//				"--runtime-set",
		//				"ignore_warnings_on_exit",
		//				"true", "src/", "lib/", "tests/",
		//			}
		//			cmdarg := CmdArg{cmd, results}
		//			So(splitCommand(command), ShouldResemble, cmdarg)
		//		})

		Convey("clean vendor bin", func() {
			So(cleanVendorBin("./vendor/bin/baz"), ShouldEqual, "./baz")
		})

		Convey("no need to clean vendor bin", func() {
			So(cleanVendorBin("./bin/baz"), ShouldEqual, "./bin/baz")
		})

		Convey("run and echo a php -v command", func() {
			//c := CmdArg{"grep", []string{"-v"}}
			//executeCommands(c)
			//So(executeCommands(c), ShouldNotPanic)
		})

		Reset(func() {
			os.Remove("/tmp/circle.yml")
		})
	})
}
