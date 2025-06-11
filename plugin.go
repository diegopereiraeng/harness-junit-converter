// plugin.go
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
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
		NestedJsonList         bool
		TestJUnitSkipField     string
		Status                 Status
	}
	Output struct {
		OutputFile string
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
		Time      int      `xml:"time,attr"`
		Name      string   `xml:"name,attr"`
		Classname string   `xml:"classname,attr"`
		Failure   *Failure `xml:"failure"`
	}
	Failure struct {
		Text    string `xml:",chardata"`
		Message string `xml:"message,attr"`
	}
)

type Plugin struct {
	Config Config
}

type Status struct {
	Total  int
	Passed int
	Errors int
	Score  float64
}

var status Status

func printHeader() {
	fmt.Println("|----------------------------------|")
	fmt.Println("|  Harness JUnit Converter Plugin  |")
	fmt.Println("|----------------------------------|")
	fmt.Println("|     Developer: Diego Pereira     |")
	fmt.Println("|----------------------------------|")
	fmt.Println("|     Version: 1.0.1               |")
	fmt.Println("|----------------------------------|")
	fmt.Println("|     Date: 2021-09-01             |")
	fmt.Println("|----------------------------------|")
	fmt.Println("")
	fmt.Println("")
}

func printStatusTable(status Status) {
	fmt.Println("|----------------------------------|")
	fmt.Println("|             Status               |")
	fmt.Println("|----------------------------------|")
	fmt.Printf("  Total:   %-3d                    \n", status.Total)
	fmt.Printf("  Passed:  %-3d                    \n", status.Passed)
	fmt.Printf("  Errors:  %-3d                    \n", status.Errors)
	fmt.Println("|----------------------------------|")
	fmt.Printf("  Score:   %-7.2f                 \n", status.Score)
	fmt.Println("|----------------------------------|")
}

func exportMetricsToFile(status Status) {
	file, err := os.Create("metrics.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()
	fmt.Fprintf(file, "TOTAL=%d\n", status.Total)
	fmt.Fprintf(file, "PASSED=%d\n", status.Passed)
	fmt.Fprintf(file, "ERRORS=%d\n", status.Errors)
	fmt.Fprintf(file, "SCORE=%.2f\n", status.Score)
}

func (p *Plugin) Exec() error {
	printHeader()

	var jsonContent string
	if p.Config.JsonFileName != "" {
		jsonRead, err := ReadJSON(p.Config.JsonFileName)
		if err != nil {
			return fmt.Errorf("error reading JSON file: %s", err)
		}
		jsonContent = jsonRead
	} else if p.Config.JsonContent != "" {
		jsonContent = p.Config.JsonContent
	} else {
		return fmt.Errorf("either JsonFileName or JsonContent must be specified")
	}

	fmt.Println("Parsing JSON to JUnit...")
	junitReport, err := ParseJunit(jsonContent, p.Config)
	if err != nil {
		return fmt.Errorf("error parsing JSON to JUnit: %s", err)
	}

	junitXML, err := xml.MarshalIndent(junitReport, " ", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JUnit to XML: %s", err)
	}

	if err := os.WriteFile(p.Config.TestName+"-junit.xml", junitXML, 0644); err != nil {
		return fmt.Errorf("error writing JUnit XML to file: %s", err)
	}

	fmt.Println("|---------------------------------------------------------------------------|")
	fmt.Println(string(junitXML))
	fmt.Println("-----------------------------------------------------------------------------")
	printStatusTable(status)

	if p.Config.FailOnFailure && junitReport.TestSuite[0].Errors > 0 {
		fmt.Println("Fail on Error Setting is True")
		return fmt.Errorf("error: There are errors in the JUnit report")
	}

	fmt.Println("Plugin executed successfully!")
	return nil
}

func ParseJunit(jsonContent string, settings Config) (*Testsuites, error) {
	failed := 0
	total := 0

	// unmarshal root and list
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonContent), &result)
	var resultList []interface{}
	json.Unmarshal([]byte(jsonContent), &resultList)

	// pick root suite name & desc
	testSuiteName, ok := result[settings.TestJUnitName].(string)
	if !ok {
		testSuiteName = settings.TestJUnitName
	}

	// new — always empty if missing or wrong type
	desc := ""
	if v, ok := result[settings.TestDescription].(string); ok {
	    desc = v
	}

	// compute suiteTime + fallback to 1
	testSuiteTime := 0
	if f, ok := result[settings.TestJUnitTime].(float64); ok {
		testSuiteTime = int(f)
	} else if s := settings.TestJUnitTime; s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			testSuiteTime = n
		}
	}
	if testSuiteTime < 1 {
		testSuiteTime = 1
	}

	// build list of suites
	var testSuiteList []interface{}
	if settings.TestJUnitList != "." && !settings.NestedJsonList {
		if arr, ok := result[settings.TestJUnitList].([]interface{}); ok {
			testSuiteList = arr
		} else {
			return nil, fmt.Errorf("failed to parse TestJUnitList as []interface{}")
		}
	} else {
		testSuiteList = resultList
	}

	testSuites := &Testsuites{}
	if len(testSuiteList) > 0 && settings.NestedJsonList {
		testSuites.TestSuite = make([]Testsuite, len(testSuiteList))
	} else {
		testSuites.TestSuite = make([]Testsuite, 1)
	}

	if settings.NestedJsonList {
		inc := 0
		for _, raw := range testSuiteList {
			total++
			m := raw.(map[string]interface{})

			// name
			name := settings.TestJUnitName
			if v, ok := m[settings.TestJUnitName].(string); ok && v != "" {
				name = v
			}

			// **package** lookup
			pkg := settings.TestJUnitPackage
			if v, ok := m[settings.TestJUnitPackage].(string); ok && v != "" {
				pkg = v
			}

			// suite-time + fallback
			st := 0
			if f, ok := m[settings.TestJUnitTime].(float64); ok {
				st = int(f)
			}
			if st < 1 {
				st = 1
			}

			testSuites.TestSuite[inc].Name = name
			testSuites.TestSuite[inc].Package = pkg
			testSuites.TestSuite[inc].Time = st

			// inner cases list
			var cases []interface{}
			if arr, ok := m[settings.TestJUnitList].([]interface{}); ok {
				cases = arr
			} else {
				cases = resultList
			}

			testSuites.TestSuite[inc].Tests = len(cases)

			for _, rc := range cases {
				// each case
				cm := rc.(map[string]interface{})
				total++

				// skip?
				if settings.TestJUnitSkipField != "" {
					if skip, ok := cm[settings.TestJUnitSkipField].(bool); ok && skip {
						continue
					}
				}

				// case name/class
				cn, _ := cm[settings.TestJUnitListName].(string)
				cc, _ := cm[settings.TestJUnitListClassName].(string)

				// **case-time + fallback**
				ct := 0
				if f, ok := cm[settings.TestJUnitListTime].(float64); ok {
					ct = int(f)
				}
				if ct < 1 {
					ct = 1
				}

				tc := Testcase{Name: cn, Classname: cc, Time: ct}

				// failure?
				if msg, ok := cm[settings.TestJUnitListFailure].(string); ok && msg != "" {
					tc.Failure = &Failure{Message: msg}
					failed++
				}

				testSuites.TestSuite[inc].TestCase = append(testSuites.TestSuite[inc].TestCase, tc)
			}

			testSuites.TestSuite[inc].Errors = failed
			inc++
		}

	} else {
		// single suite path
		testSuites.TestSuite[0].Name = testSuiteName

		// **dynamic package**
		if pkg, ok := result[settings.TestJUnitPackage].(string); ok && pkg != "" {
			testSuites.TestSuite[0].Package = pkg
		} else {
			testSuites.TestSuite[0].Package = settings.TestJUnitPackage
		}
		testSuites.TestSuite[0].Package = desc
		testSuites.TestSuite[0].Time = testSuiteTime
		testSuites.TestSuite[0].Tests = len(testSuiteList)

		for _, raw := range testSuiteList {
			total++
			cm := raw.(map[string]interface{})

			if settings.TestJUnitSkipField != "" {
				if skip, ok := cm[settings.TestJUnitSkipField].(bool); ok && skip {
					continue
				}
			}

			cn, _ := cm[settings.TestJUnitListName].(string)
			cc, _ := cm[settings.TestJUnitListClassName].(string)

			ct := 0
			if f, ok := cm[settings.TestJUnitListTime].(float64); ok {
				ct = int(f)
			}
			if ct < 1 {
				ct = 1
			}

			tc := Testcase{Name: cn, Classname: cc, Time: ct}
			if msg, ok := cm[settings.TestJUnitListFailure].(string); ok && msg != "" {
				tc.Failure = &Failure{Message: msg}
				failed++
			}

			testSuites.TestSuite[0].TestCase = append(testSuites.TestSuite[0].TestCase, tc)
		}
		testSuites.TestSuite[0].Errors = failed
	}

	status = Status{
		Total:  total,
		Passed: total - failed,
		Errors: failed,
		Score:  float64(total-failed) / float64(total) * 100,
	}

	return testSuites, nil
}

func ReadJSON(filename string) (string, error) {
	result, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
