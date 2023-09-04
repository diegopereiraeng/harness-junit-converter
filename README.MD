
# Harness JUnit Converter Plugin

This plugin is designed to convert JSON fields to the JUnit XML format, suitable for use within a Harness CI/CD pipeline. The plugin is built using Go and is packaged as a Docker container for effortless integration.

## Parameters

The tool accepts the following parameters:

- **json_file_name**: Name of the JSON file.
- **json_content**: Direct JSON content.
- **test_name**: Name of the test.
- **test_description**: Description of the test (optional).
- **test_junit_time**: JUnit time.
- **test_junit_package**: JUnit package name.
- **test_junit_name**: JUnit name.
- **test_junit_list**: Json Path to the test list.
- **test_junit_list_name**: Name for JUnit list.
- **test_junit_list_class_name**: Class name for JUnit list.
- **test_junit_list_failure**: Failure message for JUnit list.
- **test_junit_list_time**: Time for each JUnit list test.

## Additional Parameters 

- **test_junit_skip_field**: Json path to the field that would give true or false to skip errors.
- **nested_json_list**: (true|false) If you have a nested json that is a list of a Test Suite, than use this option

## JSON List Support

e.g: [{"name": "value", "desc": "test2",...},{...}]

- **test_junit_list**: use this one with value "."

``` yaml
test_junit_list = "."
```

## Nested JSon example

Look the file in the repo tests/nested-json.json

``` yaml
nested_json_list = true
```

### Kube Score JSON Convert Example

``` yaml
step:
  name: Kube Score Ingestion
  identifier: kube_score
  type: Plugin
  spec:
    connectorRef: <Docker_Registry>
    image: diegokoala/harness-junit-converter:latest
    reports:
      type: JUnit
      spec:
        paths:
          - "**/**/*.xml"
    settings:
      json_file_name: 'tests/nested-json.json'
      test_name: "object_name"
      test_description: "file_name"
      test_junit_time: "file_row"
      test_junit_package: "file_name"
      test_junit_name: "object_name"
      test_junit_list: "checks"
      test_junit_list_name: "name"
      test_junit_list_class_name: "comment"
      test_junit_list_failure: "comments[].summary"
      test_junit_list_time: "grade"
      nested_json_list: true
      test_junit_skip_field: "skipped"
```

## How to Use This Plugin in a Harness Pipeline

### YAML Configuration Example

Here is a full example of how to configure the plugin within a Harness pipeline:

```yaml
type: Plugin
spec:
    connectorRef: <Docker_Hub_Account_or_Other_Registry>
    image: diegokoala/harness-junit-converter:latest
    reports:
        type: JUnit
        spec:
            paths:
              - "**/**/*.xml"
    settings:
        json_file_name: <Path_to_JSON_File>
        json_content: <JSON_Content_As_String>
        test_name: <JSON_Field_For_Test_Name>
        test_description: <JSON_Field_For_Test_Description>
        test_junit_time: <JSON_Field_For_JUnit_Time>
        test_junit_package: <JSON_Field_For_JUnit_Package>
        test_junit_name: <JSON_Field_For_JUnit_Name>
        test_junit_list: <JSON_Field_For_JUnit_Test_List>
        test_junit_list_name: <JSON_Field_For_Test_Case_Name>
        test_junit_list_class_name: <JSON_Field_For_Test_Case_Class_Name>
        test_junit_list_failure: <JSON_Field_For_Test_Case_Failure_Message>
        test_junit_list_time: <JSON_Field_For_Test_Case_Time>
```

### Parameters Explained

- `connectorRef`: Your Docker Hub account or other container registry.
- `image`: The Docker image of the plugin, usually hosted on Docker Hub.
- `json_file_name`: The path to the JSON file that will be converted to JUnit format. This is optional if `json_content` is provided.
- `json_content`: The JSON content as a string that will be converted to JUnit format. This is optional if `json_file_name` is provided.
- Other `settings`: These act as JSON paths that map your JSON fields to JUnit fields. For example, if your JSON has a field `myField` that you want to map to `test_name`, you would set `test_name: myField`.

### Additional Example

Here is an additional example to provide more context on how to use the plugin:

```yaml
type: Plugin
spec:
    connectorRef: MyDockerHubAccount
    image: diegokoala/harness-junit-converter:latest
    reports:
        type: JUnit
        spec:
            paths:
              - "target/**/test-*.xml"
    settings:
        json_content: '{"test_suite": "MySuite", "packageOrOtherThing": "value" "test_list": [{"name": "Test1", "class": "Class1", "time": 5}, {"name": "Test2", "class": "Class2", "failure": "failure msg"}]}'
        test_name: "Custom Validation"
        test_description: "Validate errors in the custom JSON"
        test_junit_list: "test_list"
        test_junit_list_name: "name"
        test_junit_list_class_name: "class"
        test_junit_list_time: "time"
        test_junit_package: "packageOrOtherThing"
        test_junit_name: "test_suite"
        test_junit_list_failure: "test_junit_list_failure"
```

### Another Additional Example

Here is an another additional example to provide more context on how to use the plugin:
```
{
  "test_junit_name": "MyTestSuite",
  "test_junit_package": "com.example",
  "test_junit_time": 50,
  "test_junit_list": [
    {
      "test_junit_list_name": "TestCase1",
      "test_junit_list_class_name": "TestClass1",
      "test_junit_list_failure": "",
      "test_junit_list_time": 10
    },
    {
      "test_junit_list_name": "TestCase2",
      "test_junit_list_class_name": "TestClass2",
      "test_junit_list_failure": "Some Failure",
      "test_junit_list_time": 15
    },
    {
      "test_junit_list_name": "TestCase3",
      "test_junit_list_class_name": "TestClass3",
      "test_junit_list_failure": "",
      "test_junit_list_time": 25
    }
  ]
}
```

### Example using the Json Above

```yaml
type: Plugin
spec:
    connectorRef: <+input>
    image: diegokoala/harness-junit-converter:latest
    reports:
        type: JUnit
        spec:
            paths:
              - "**/**/*.xml"
    settings:
        json_content: '{"test_junit_name": "MyTestSuite", "test_junit_package": "com.example", "test_junit_time": 50, "test_junit_list": [{"test_junit_list_name": "TestCase1", "test_junit_list_class_name": "TestClass1", "test_junit_list_failure": "", "test_junit_list_time": 10}, {"test_junit_list_name": "TestCase2", "test_junit_list_class_name": "TestClass2", "test_junit_list_failure": "Some Failure", "test_junit_list_time": 15}, {"test_junit_list_name": "TestCase3", "test_junit_list_class_name": "TestClass3", "test_junit_list_failure": "", "test_junit_list_time": 25}]}'
        test_name: "test_junit_name"
        test_description: "test_junit_package"
        test_junit_time: "test_junit_time"
        test_junit_package: "test_junit_package"
        test_junit_name: "test_junit_name"
        test_junit_list: "test_junit_list"
        test_junit_list_name: "test_junit_list_name"
        test_junit_list_class_name: "test_junit_list_class_name"
        test_junit_list_failure: "test_junit_list_failure"
        test_junit_list_time: "test_junit_list_time"
```



### Running the Plugin Locally

To test the plugin locally, you can build and run it using Go and Docker:

```bash
go mod download
go build -o harness-junit-converter
./harness-junit-converter --json_content='<Your_JSON_Content_Here>'
```

To build and run the Docker container:

```bash
docker build -t diegokoala/harness-junit-converter:latest .
docker run -e PLUGIN_JSON_CONTENT='{"test_junit_name": "MyTestSuite", "test_junit_package": "com.example", "test_junit_time": 50, "test_junit_list": [{"test_junit_list_name": "TestCase1", "test_junit_list_class_name": "TestClass1", "test_junit_list_failure": "", "test_junit_list_time": 10}, {"test_junit_list_name": "TestCase2", "test_junit_list_class_name": "TestClass2", "test_junit_list_failure": "Some Failure", "test_junit_list_time": 15}, {"test_junit_list_name": "TestCase3", "test_junit_list_class_name": "TestClass3", "test_junit_list_failure": "", "test_junit_list_time": 25}]}' \
           -e PLUGIN_TEST_NAME='test_junit_name' \
           -e PLUGIN_TEST_DESCRIPTION='test_junit_package' \
           -e PLUGIN_TEST_JUNIT_TIME='test_junit_time' \
           -e PLUGIN_TEST_JUNIT_PACKAGE='test_junit_package' \
           -e PLUGIN_TEST_JUNIT_NAME='test_junit_name' \
           -e PLUGIN_TEST_JUNIT_LIST='test_junit_list' \
           -e PLUGIN_TEST_JUNIT_LIST_NAME='test_junit_list_name' \
           -e PLUGIN_TEST_JUNIT_LIST_CLASS_NAME='test_junit_list_class_name' \
           -e PLUGIN_TEST_JUNIT_LIST_FAILURE='test_junit_list_failure' \
           -e PLUGIN_TEST_JUNIT_LIST_TIME='test_junit_list_time' \
           diegokoala/harness-junit-converter:latest

```

## License

This project is licensed under the MIT License.

## Contributing

Pull requests are welcome. For significant changes, please open an issue first to discuss what you would like to change.
