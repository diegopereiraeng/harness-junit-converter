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
//
// Config struct holds the configuration options for the plugin.
// TestName: the name of the test.
// TestDescription: the description of the test.
// TestJUnitTime: the time taken by the test.
// TestJUnitPackage: the package of the test.
// TestJUnitName: the name of the test.
// TestJUnitList: the list of tests.
// TestJUnitListName: the name of the test list.
// TestJUnitListClassName: the class name of the test list.
// TestJUnitListFailure: the failure of the test list.
// TestJUnitListTime: the time taken by the test list.
// JsonFileName: the name of the JSON file.
// JsonContent: the content of the JSON file.
// FailOnFailure: whether to fail on failure.
// NestedJsonList: whether the JSON list is nested.
// TestJUnitSkipField: the field to skip in the JUnit report.
//
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
	"os"
	"strconv"
	"strings"
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
		NestedJsonList         bool
		TestJUnitSkipField     string
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
		return fmt.Errorf("either JsonFileName or JsonContent must be specified")
	}

	jsonContent := ""

	// code that will convert the JSON to JUnit
	if p.Config.JsonFileName != "" {
		// Read the JSON file
		jsonRead, err := ReadJSON(p.Config.JsonFileName)
		if err != nil {
			return fmt.Errorf("error reading JSON file: %s", err)
		}
		jsonContent = jsonRead

	} else if p.Config.JsonContent != "" {
		// Use the direct JSON content
		jsonContent = p.Config.JsonContent
	} else {
		return fmt.Errorf("either JsonFileName or JsonContent must be specified")
	}

	// Parse JSON to JUnit
	fmt.Println("Parsing JSON to JUnit...")
	junitReport, err := ParseJunit(jsonContent, p.Config)
	if err != nil {
		return fmt.Errorf("error parsing JSON to JUnit: %s", err)
	}

	// Serialize JUnit to XML and print (or write to file)
	junitXML, err := xml.MarshalIndent(junitReport, " ", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JUnit to XML: %s", err)
	}
	fmt.Println(string(junitXML))

	//save to a file called <test_name_variable>-junit.xml
	err = os.WriteFile(p.Config.TestName+"-junit.xml", junitXML, 0644)
	if err != nil {
		return fmt.Errorf("error writing JUnit XML to file: %s", err)
	}

	// Print the plugin config

	fmt.Println("Plugin executed with config:", p.Config)

	// Check if should fail on errors
	if p.Config.FailOnFailure {
		// verify if there are errors in Testsuites object
		if junitReport.TestSuite[0].Errors > 0 {
			fmt.Println("Fail on Error Setting is True")
			fmt.Println("Error: There are errors in the JUnit report.")
			return fmt.Errorf("error: There are errors in the JUnit report")
		}

	}

	// Verify that the plugin works
	fmt.Println("Plugin executed successfully!")

	return nil
}

func ParseJunit(jsonContent string, settings Config) (*Testsuites, error) {

	// Add support to this json:
	// [{"code":"DL3018","column":1,"file":"Dockerfile","level":"warning","line":4,"message":"Pin versions in apk add. Instead of `apk add <package>` use `apk add <package>=<version>`"},{"code":"DL3059","column":1,"file":"Dockerfile","level":"info","line":17,"message":"Multiple consecutive `RUN` instructions. Consider consolidation."}]

	// JUnit conversion logic
	failed := 0
	total := 0 // Count for total test cases
	testCases := []Testcase{}
	errors := 0
	newError := 0

	// Parse the JSON content
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonContent), &result)
	var resultList []interface{}
	json.Unmarshal([]byte(jsonContent), &resultList)
	// fmt.Println("resultList: ", resultList)

	// Get the test suite name
	testSuiteName, ok := result[settings.TestJUnitName].(string)
	if !ok {
		testSuiteName = settings.TestJUnitName
	}

	// Get the test suite description
	testSuiteDescription, ok := result[settings.TestDescription].(string)
	if !ok {
		testSuiteDescription = settings.TestDescription
	}

	// Get the test suite time
	testSuiteTime := 0
	testSuiteTimeFloat, ok := result[settings.TestJUnitTime].(float64)
	if ok {
		testSuiteTime = int(testSuiteTimeFloat)
	} else {
		fmt.Println("TestJUnitTime is not float64")
		fmt.Println("TestJUnitTime: ", settings.TestJUnitTime)
		testSuiteTimeInt, ok := result[settings.TestJUnitTime].(int)
		if ok {
			testSuiteTime = testSuiteTimeInt
		} else {
			if settings.TestJUnitTime != "" {
				fmt.Println("TestJUnitTime is not empty")
				fmt.Println("TestJUnitTime: ", settings.TestJUnitTime)
				testSuiteTimeInt, err := strconv.Atoi(settings.TestJUnitTime)
				if err != nil {
					// return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
					fmt.Println("Error: failed to parse TestJUnitTime as float64 or int")
					fmt.Println("Error: ", err)
					fmt.Println("setting a default value of 0")
					testSuiteTimeInt = 0
				}
				testSuiteTime = testSuiteTimeInt
			} else {
				// return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
				fmt.Println("Error: failed to parse TestJUnitTime as float64 or int")
				fmt.Println("setting a default value of 0")
				testSuiteTime = 0

			}
		}
	}

	// Get the test suite list
	var testSuiteList []interface{}
	if settings.TestJUnitList != "." && !settings.NestedJsonList {
		fmt.Println("TestJUnitList is not . and NestedJsonList is false")
		testSuiteListInterface, ok := result[settings.TestJUnitList].([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse TestJUnitList as []interface{}")
		}
		fmt.Println("TestSuiteListInterface: ", testSuiteListInterface)
		testSuiteList = testSuiteListInterface
	} else {
		fmt.Println("TestJUnitList is . or NestedJsonList is true")
		testSuiteList = resultList
		fmt.Println("Assigning resultList to testSuiteList")
	}

	// Create the testsuites object
	testSuites := &Testsuites{}

	// Create the
	fmt.Println("len(testSuiteList): ", len(testSuiteList))
	if len(testSuiteList) > 0 && settings.NestedJsonList {
		fmt.Println("len(testSuiteList) > 0 and NestedJsonList is true")
		testSuites.TestSuite = make([]Testsuite, len(testSuiteList))
	} else {
		fmt.Println("len(testSuiteList) <= 0 or NestedJsonList is false")
		testSuites.TestSuite = make([]Testsuite, 1)
		fmt.Println("len(testSuites.TestSuite): ", len(testSuites.TestSuite))
	}
	// testSuites.TestSuite = make([]Testsuite, len(testSuiteList)+1)
	fmt.Println("NestedJsonList: ", settings.NestedJsonList)
	fmt.Println("len(testSuites.TestSuite): ", len(testSuites.TestSuite))
	if settings.NestedJsonList {
		fmt.Println("NestedJsonList is true")
		inc := 0
		// Iterate over the test suites
		for _, testSuite := range testSuiteList {
			// fmt.Println("TestSuite: ", testSuite)
			total++ // Increment the total test cases count

			testCaseMap := testSuite.(map[string]interface{})

			testSuiteName, ok := testCaseMap[settings.TestJUnitName].(string)
			if !ok || testSuiteName == "" {
				testSuiteName = settings.TestJUnitName
			}
			if !ok {
				testSuiteName = settings.TestJUnitName
			}
			fmt.Println("Name: ", testSuiteName)
			testSuites.TestSuite[inc].Name = testSuiteName

			testSuiteDescription, ok := testCaseMap[settings.TestDescription].(string)
			if !ok {
				testSuiteDescription = settings.TestDescription
			}
			fmt.Println("Description: ", testSuiteDescription)
			testSuites.TestSuite[inc].Package = testSuiteDescription

			testSuiteTime := 0
			testSuiteTimeFloat, ok := testCaseMap[settings.TestJUnitTime].(float64)
			if ok {
				testSuiteTime = int(testSuiteTimeFloat)
			} else {
				fmt.Println("TestJUnitTime is not float64")
				fmt.Println("TestJUnitTime: ", settings.TestJUnitTime)
				testSuiteTimeInt, ok := testCaseMap[settings.TestJUnitTime].(int)
				if ok {
					testSuiteTime = testSuiteTimeInt
				} else {
					if settings.TestJUnitTime != "" {
						fmt.Println("TestJUnitTime is not empty")
						fmt.Println("TestJUnitTime: ", settings.TestJUnitTime)
						testSuiteTimeInt, err := strconv.Atoi(settings.TestJUnitTime)
						if err != nil {
							// return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
							fmt.Println("Error: failed to parse TestJUnitTime as float64 or int")
							fmt.Println("Error: ", err)
							fmt.Println("setting a default value of 0")
							testSuiteTimeInt = 0
						}
						testSuiteTime = testSuiteTimeInt
					} else {
						// return nil, fmt.Errorf("failed to parse TestJUnitTime as float64 or int")
						fmt.Println("Error: failed to parse TestJUnitTime as float64 or int")
						fmt.Println("setting a default value of 0")
						testSuiteTime = 0

					}
				}
			}
			fmt.Println("Time: ", testSuiteTime)
			testSuites.TestSuite[inc].Time = testSuiteTime

			// Get the test suite list
			var testSuiteList []interface{}
			if settings.TestJUnitList != "." && !settings.NestedJsonList {
				fmt.Println("TestJUnitList is not . and NestedJsonList is false")
				testSuiteListInterface, ok := testCaseMap[settings.TestJUnitList].([]interface{})
				if !ok {
					return nil, fmt.Errorf("failed to parse TestJUnitList as []interface{}")
				}
				fmt.Println("TestSuiteListInterface: ", testSuiteListInterface)
				testSuiteList = testSuiteListInterface
			} else {
				fmt.Println("TestJUnitList is . or NestedJsonList is true")
				testSuiteList = resultList
			}

			testSuites.TestSuite[inc].Tests = len(testSuiteList)

			// Iterate over the test cases

			testCases := []Testcase{}

			// Populate fields for a single test suite
			singleTestSuite := Testsuite{
				Name:    testSuiteName,
				Package: testSuiteDescription,
				Time:    testSuiteTime,
				Tests:   len(testSuiteList),
			}

			// Add test cases (assuming they are in a list for each object)
			testCaseList := testCaseMap[settings.TestJUnitList].([]interface{})
			for _, testCase := range testCaseList {
				fmt.Println("Case: ", testCase)

				testCaseMap := testCase.(map[string]interface{})

				// Remove "s" in the end of settings.TestJUnitList
				nestedJson := settings.TestJUnitList[:len(settings.TestJUnitList)-1]
				fmt.Println("nestedJson: ", nestedJson)
				var testCaseObj Testcase
				// var testCaseName, testCaseClassName string
				// var testCaseTime int
				// var checkMapMap map[string]interface{}

				var time int

				if tempTime, ok := testCaseMap[settings.TestJUnitListTime]; ok && tempTime != nil {
					if timeFloat, ok := tempTime.(float64); ok {
						time = int(timeFloat)
					} else {
						return nil, fmt.Errorf("failed to cast TestJUnitListTime to float64")
					}
				} else {
					fmt.Println("settings.TestJUnitListTime: ", settings.TestJUnitListTime)
					fmt.Println("Convert to int")
					// var err error
					// conver checkMap[settings.TestJUnitListTime] to int
					fmt.Println("settings.TestJUnitListTime: ", settings.TestJUnitListTime)
					fmt.Println("checkMap[settings.TestJUnitListTime]: ", testCaseMap[settings.TestJUnitListTime])
					time = testCaseMap[settings.TestJUnitListTime].(int)

					// time, err = checkMap[settings.TestJUnitListTime].(int)
					// if err != nil {
					// 	return nil, fmt.Errorf("failed to parse TestJUnitListTime as float64 or int")
					// }
					// return nil, fmt.Errorf("TestJUnitListTime is either nil or not found in checkMap")
				}

				if checkMap, ok := testCaseMap[nestedJson].(map[string]interface{}); ok {

					// check if checkMap["map"] is null or empty
					// if it is, then checkMap["map"] = checkMap
					// if it is not, then checkMap["map"] = checkMap["map"]
					// checkMap["map"] = checkMap
					// checkMap["map"] = checkMap["map"]
					fmt.Println("checkMap: ", checkMap)
					if checkMap[settings.TestJUnitListName] == nil {
						// checkMap["map"] = checkMap
						//skip the loop
						fmt.Println("checkMap[" + settings.TestJUnitListName + "] is nil")
						fmt.Println("Continuing...")
						continue

					}

					// checkMapMap = checkMap["map"].(map[string]interface{})

					// fmt.Println("checkMapMap: ", checkMapMap)
					// fmt.Println("settings: " + settings.TestJUnitListName)
					// fmt.Println("checkMapMap[settings.TestJUnitListName]: ", checkMapMap[settings.TestJUnitListName])
					fmt.Println("checkMap[settings.TestJUnitListName]: ", checkMap[settings.TestJUnitListName])
					var name, classname string

					if tempName, ok := checkMap[settings.TestJUnitListName]; ok && tempName != nil {
						if nameStr, ok := tempName.(string); ok {
							name = nameStr
						} else {
							return nil, fmt.Errorf("failed to cast TestJUnitListName to string")
						}
					} else {
						return nil, fmt.Errorf("TestJUnitListName is either nil or not found in checkMap")
					}

					if tempClass, ok := checkMap[settings.TestJUnitListClassName]; ok && tempClass != nil {
						if classStr, ok := tempClass.(string); ok {
							classname = classStr
						} else {
							return nil, fmt.Errorf("failed to cast TestJUnitListClassName to string")
						}
					} else {
						return nil, fmt.Errorf("TestJUnitListClassName is either nil or not found in checkMap")
					}

					testCaseObj = Testcase{
						Name:      name,
						Classname: classname,
						Time:      time,
					}

				}
				var failureMessage string
				fmt.Println("settings.TestJUnitListFailure: ", settings.TestJUnitListFailure)
				var skipFieldList bool
				if settings.TestJUnitSkipField != "" {
					skipFieldList = testCaseMap[settings.TestJUnitSkipField].(bool)
				} else {
					skipFieldList = false
				}
				error := 0
				if strings.Contains(settings.TestJUnitListFailure, "[]") && !skipFieldList {
					fmt.Println("Failures are in a list")

					// This is a list, handle accordingly
					failureMessageList := []string{}
					// split the string by [] and get the first part
					nestedList := strings.Split(settings.TestJUnitListFailure, "[]")
					nestedListFirst := nestedList[0]
					splitedNestedListSecond := strings.Split(nestedList[1], ".")
					nestedListSecond := splitedNestedListSecond[1]
					fmt.Println("nestedListFirst: ", nestedListFirst)
					fmt.Println("nestedListSecond: ", nestedListSecond)

					errs := 0
					if commentsList, ok := testCaseMap[nestedListFirst].([]interface{}); ok {
						fmt.Println("commentsList: ", commentsList)
						if len(commentsList) > 0 {
							errs = len(commentsList)
						}
						for _, comment := range commentsList {
							if commentMap, ok := comment.(map[string]interface{}); ok {
								if summary, ok := commentMap[nestedListSecond].(string); ok {
									failureMessageList = append(failureMessageList, summary)
								}
							}
						}
					}
					if errs > 0 {
						failureMessage = strings.Join(failureMessageList, "; ")
						error++
					}
				} else {

					if !skipFieldList {
						// Existing single value handling
						if testCaseMap[settings.TestJUnitListFailure] != nil {
							failureMessage = testCaseMap[settings.TestJUnitListFailure].(string)
							error++
						}
						// else {
						// 	failureMessage = "Unknown Error"
						// 	fmt.Println("Error: TestJUnitListFailure is nil")
						// 	fmt.Println("Error: ", testCaseMap[settings.TestJUnitListFailure])
						// 	fmt.Println("Error: ", testCaseMap[settings.TestJUnitListFailure].(string))

						// }
					}
				}
				if !skipFieldList && error > 0 {
					failed++
					errors++
					testCaseObj.Failure = &Failure{Message: failureMessage}
				}

				testCases = append(testCases, testCaseObj)

			}

			singleTestSuite.Errors = errors
			singleTestSuite.TestCase = testCases

			// Append this singleTestSuite to the main testSuites object.
			testSuites.TestSuite = append(testSuites.TestSuite, singleTestSuite)
			inc++
		}
	} else {
		fmt.Println("NestedJsonList is false")
		testSuites.TestSuite[0].Name = testSuiteName
		testSuites.TestSuite[0].Package = testSuiteDescription
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
	}

	return testSuites, nil
}

func ReadJSON(filename string) (string, error) {

	// Read the JSON file and return its contents as a string
	result, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
