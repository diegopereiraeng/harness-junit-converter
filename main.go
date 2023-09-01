package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var build = "1" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "Harness-JUnit-Converter"
	app.Usage = "CLI tool to convert JSON fields to JUnit XML format."
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "json_file_name",
			Usage:  "Name of the JSON file.",
			EnvVar: "PLUGIN_JSON_FILE_NAME",
		},
		cli.StringFlag{
			Name:   "json_content",
			Usage:  "Direct JSON content.",
			EnvVar: "PLUGIN_JSON_CONTENT",
		},
		cli.StringFlag{
			Name:   "test_name",
			Usage:  "Name of the test.",
			EnvVar: "PLUGIN_TEST_NAME",
		},
		cli.StringFlag{
			Name:   "test_description",
			Usage:  "Description of the test (optional).",
			EnvVar: "PLUGIN_TEST_DESCRIPTION",
		},
		cli.StringFlag{
			Name:   "test_junit_time",
			Usage:  "JUnit time.",
			EnvVar: "PLUGIN_TEST_JUNIT_TIME",
		},
		cli.StringFlag{
			Name:   "test_junit_package",
			Usage:  "JUnit package name.",
			EnvVar: "PLUGIN_TEST_JUNIT_PACKAGE",
		},
		cli.StringFlag{
			Name:   "test_junit_name",
			Usage:  "JUnit name.",
			EnvVar: "PLUGIN_TEST_JUNIT_NAME",
		},
		cli.StringSliceFlag{
			Name:   "test_junit_list",
			Usage:  "List of JUnit tests.",
			EnvVar: "PLUGIN_TEST_JUNIT_LIST",
		},
		cli.StringFlag{
			Name:   "test_junit_list_name",
			Usage:  "Name for JUnit list.",
			EnvVar: "PLUGIN_TEST_JUNIT_LIST_NAME",
		},
		cli.StringFlag{
			Name:   "test_junit_list_class_name",
			Usage:  "Class name for JUnit list.",
			EnvVar: "PLUGIN_TEST_JUNIT_LIST_CLASS_NAME",
		},
		cli.StringFlag{
			Name:   "test_junit_list_failure",
			Usage:  "Failure message for JUnit list.",
			EnvVar: "PLUGIN_TEST_JUNIT_LIST_FAILURE",
		},
		cli.StringFlag{
			Name:   "test_junit_list_time",
			Usage:  "Time for each JUnit list test.",
			EnvVar: "PLUGIN_TEST_JUNIT_LIST_TIME",
		},
		cli.BoolFlag{
			Name:   "fail_on_errors",
			Usage:  "Fail the execution on errors.",
			EnvVar: "PLUGIN_FAIL_ON_ERRORS",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	if c.String("json_file_name") != "" && c.String("json_content") != "" {
		fmt.Println("Error: Please specify either json_file_name or json_content, but not both.")
		os.Exit(1)
	}

	config := Config{
		TestName:               c.String("test_name"),
		TestDescription:        c.String("test_description"),
		TestJUnitTime:          c.String("test_junit_time"),
		TestJUnitPackage:       c.String("test_junit_package"),
		TestJUnitName:          c.String("test_junit_name"),
		TestJUnitList:          c.String("test_junit_list"),
		TestJUnitListName:      c.String("test_junit_list_name"),
		TestJUnitListClassName: c.String("test_junit_list_class_name"),
		TestJUnitListFailure:   c.String("test_junit_list_failure"),
		TestJUnitListTime:      c.String("test_junit_list_time"),
		JsonFileName:           c.String("json_file_name"),
		JsonContent:            c.String("json_content"),
		FailOnFailure:          c.Bool("fail_on_errors"),
	}

	plugin := Plugin{Config: config}
	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
