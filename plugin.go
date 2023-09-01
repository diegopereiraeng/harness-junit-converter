package main

// This plugin is responsible for converting JSON fields to JUnit XML format.
// The JSON fields are:
// - test_name
// - test_description
// - test_junit_time
// - test_junit_package
// - test_junit_name
// - test_junit_list
// - test_junit_list_name
// - test_junit_list_class_name
// - test_junit_list_failure
// - test_junit_list_time
//
// The JUnit XML format is:
// <testsuites>
//   <testsuite name="..." package="..." time="..." tests="..." errors="...">
//     <testcase name="..." classname="...">
//       <failure message="..."></failure>
//     </testcase>
//   </testsuite>
// </testsuites>
//
// The plugin will be invoked as follows:
// $ ./plugin --json_file_name=sample.json
// $ ./plugin --json_content='{"test_name": "test", "test_description": "test description", "test_junit_time": "1", "test_junit_package": "test package", "test_junit_name": "test name", "test_junit_list": [{"test_junit_list_name": "test list name", "test_junit_list_class_name": "test list class name", "test_junit_list_failure": "test list failure", "test_junit_list_time": "1"}]}'
// $ ./plugin --json_file_name=sample.json --test_name=test --test_description="test description"

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
)

type (
	Config struct {
		TestName               string
		TestDescription        string
		TestJUnitTime          string
		TestJUnitPackage       string
		TestJUnitName          string
		TestJUnitList          string
		TestJUnitListName      string
		TestJUnitListClassName string
		TestJUnitListFailure   string
		TestJUnitListTime      string
		JsonFileName           string
		JsonContent            string
		FailOnFailure          bool
	}
	Output struct {
		OutputFile string // File where plugin output are saved
	}
	Testsuites struct {
		XMLName   xml.Name    `xml:"testsuites"`
		Text      string      `xml:",chardata"`
		TestSuite []Testsuite `xml:"testsuite"`
	}
	Testsuite struct {
		Text     string     `xml:",chardata"`
		Package  string     `xml:"package,attr"`
		Time     int        `xml:"time,attr"`
		Tests    int        `xml:"tests,attr"`
		Errors   int        `xml:"errors,attr"`
		Name     string     `xml:"name,attr"`
		TestCase []Testcase `xml:"testcase"`
	}
	Testcase struct {
		Text      string   `xml:",chardata"`
		Time      int      `xml:"time,attr"`      // Actual Value Sonar
		Name      string   `xml:"name,attr"`      // Metric Key
		Classname string   `xml:"classname,attr"` // The metric Rule
		Failure   *Failure `xml:"failure"`        // Sonar Failure - show results
	}
	Failure struct {
		Text    string `xml:",chardata"`
		Message string `xml:"message,attr"`
	}
)

type Plugin struct {
	Config Config
}

func (p *Plugin) Exec() error {
	// The logic here would be similar to what was previously in the main's run function
	// Read JSON, Convert to JUnit, and Export XML
	if p.Config.JsonFileName != "" {
		// Read the JSON file
	} else if p.Config.JsonContent != "" {
		// Use the direct JSON content
	} else {
		return fmt.Errorf("Either JsonFileName or JsonContent must be specified.")
	}

	jsonContent := ""

	// code that will convert the JSON to JUnit
	if p.Config.JsonFileName != "" {
		// Read the JSON file
		jsonRead, err := ReadJSON(p.Config.JsonFileName)
		if err != nil {
			return fmt.Errorf("Error reading JSON file: %s", err)
		}
		jsonContent = jsonRead

	} else if p.Config.JsonContent != "" {
		// Use the direct JSON content
		jsonContent = p.Config.JsonContent
	} else {
		return fmt.Errorf("Either JsonFileName or JsonContent must be specified.")
	}

	// Parse JSON to JUnit
	junitReport, err := ParseJunit(jsonContent, p.Config)

	// Serialize JUnit to XML and print (or write to file)
	junitXML, err := xml.MarshalIndent(junitReport, " ", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling JUnit to XML: %s", err)
	}
	fmt.Println(string(junitXML))

	//save to a file called <test_name_variable>-junit.xml
	err = ioutil.WriteFile(p.Config.TestName+"-junit.xml", junitXML, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JUnit XML to file: %s", err)
	}

	// Print the plugin config

	fmt.Println("Plugin executed with config:", p.Config)

	// Check if should fail on errors
	if p.Config.FailOnFailure {
		// verify if there are errors in Testsuites object
		if junitReport.TestSuite[0].Errors > 0 {
			fmt.Println("Fail on Error Setting is True")
			fmt.Println("Error: There are errors in the JUnit report.")
			return fmt.Errorf("Error: There are errors in the JUnit report.")
		}

	}

	// Verify that the plugin works
	fmt.Println("Plugin executed successfully!")

	return nil
}

func ParseJunit(jsonContent string, settings Config) (*Testsuites, error) {

	// JUnit conversion logic
	failed := 0
	total := 0 // Count for total test cases
	testCases := []Testcase{}
	errors := 0
	newError := 0

	// Parse the JSON content
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonContent), &result)

	// Get the test suite name
	// testSuiteName := result["test_junit_name"].(string)
	testSuiteName, ok := result[settings.TestJUnitName].(string)
	if !ok {
		testSuiteName = settings.TestJUnitName
	}

	// Get the test suite description
	// testSuiteDescription := result["test_junit_package"].(string)
	testSuiteDescription, ok := result[settings.TestDescription].(string)
	if !ok {
		testSuiteDescription = settings.TestDescription
	}

	// Get the test suite time
	// testSuiteTime := int(result["test_junit_time"].(float64))

	// Get the test suite time
	testSuiteTime := 0
	testSuiteTimeFloat, ok := result[settings.TestJUnitTime].(float64)
	if ok {
		testSuiteTime = int(testSuiteTimeFloat)
	} else {
		testSuiteTimeInt, ok := result[settings.TestJUnitTime].(int)
		if ok {
			testSuiteTime = testSuiteTimeInt
		} else {
			if settings.TestJUnitTime != "" {
				testSuiteTimeInt, err := strconv.Atoi(settings.TestJUnitTime)
				if err != nil {
					return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
				}
				testSuiteTime = testSuiteTimeInt
			} else {
				return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
			}
		}
	}

	// Get the test suite list
	// testSuiteList := result["test_junit_list"].([]interface{})
	testSuiteList := result[settings.TestJUnitList].([]interface{})

	// Create the testsuites object
	testSuites := &Testsuites{}
	testSuites.TestSuite = make([]Testsuite, 1)
	testSuites.TestSuite[0].Name = testSuiteName
	testSuites.TestSuite[0].Package = testSuiteDescription
	//make the next optional
	testSuites.TestSuite[0].Time = testSuiteTime
	testSuites.TestSuite[0].Tests = len(testSuiteList)

	// Iterate over the test cases
	for _, testCase := range testSuiteList {
		total++ // Increment the total test cases count

		testCaseMap := testCase.(map[string]interface{})
		testCaseName := testCaseMap[settings.TestJUnitListName].(string)
		testCaseClassName := testCaseMap[settings.TestJUnitListClassName].(string)
		testCaseFailure := testCaseMap[settings.TestJUnitListFailure].(string)
		testCaseTime := int(testCaseMap[settings.TestJUnitListTime].(float64))

		// Create the testcase object
		testCaseObj := Testcase{
			Name:      testCaseName,
			Classname: testCaseClassName,
			Time:      testCaseTime,
		}

		// Check if the test case failed
		if testCaseFailure != "" {
			failed++ // Increment the failed test cases count
			errors++
			newError++
			testCaseObj.Failure = &Failure{Message: testCaseFailure}
		}

		// Add the testcase to the testsuite
		testCases = append(testCases, testCaseObj)
	}

	// Update the testsuite object
	testSuites.TestSuite[0].Errors = errors
	testSuites.TestSuite[0].TestCase = testCases

	// Print or return the total and failed test cases count
	fmt.Printf("Total test cases: %d\n", total)
	fmt.Printf("Failed test cases: %d\n", failed)

	return testSuites, nil
}

func ReadJSON(filename string) (string, error) {

	// Read the JSON file and return its contents as a string
	result, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
