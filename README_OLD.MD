# Harness JUnit Converter Plugin

This plugin is designed to convert JSON fields to JUnit XML format in a Harness pipeline. It's built in Go and packaged as a Docker container for easy integration.

## How to Use This Plugin in a Harness Pipeline

### Full Config Step Example

```yaml
type: Plugin
spec:
    connectorRef: <+input>
    image: your_docker_hub_username/harness-junit-converter:latest
    reports:
        type: JUnit
        spec:
            paths:
              - "**/**/*.xml"
    settings:
        json_file_name: <+input>
        json_content: <+input>
        test_name: <+input>
        test_description: <+input>
        test_junit_time: <+input>
        test_junit_package: <+input>
        test_junit_name: <+input>
        test_junit_list: <+input>
        test_junit_list_name: <+input>
        test_junit_list_class_name: <+input>
        test_junit_list_failure: <+input>
        test_junit_list_time: <+input>
```

The parameters in the `settings` section of the YAML file, except for `json_content` and `json_file_name`, act as JSON paths to map your JSON fields to JUnit fields. 

### Example:
Suppose you have the following JSON content:

```json
{
  "field1": "field1",
  "field2": "field2",
  "field3": "field3",
  "testList": [
    {
      "name": "test1",
      "class": "whatIWant",
      "time": 7
    },
    {
      "name": "test2",
      "class": "testClass",
      "failure": "error msg"
    }
  ]
}
```

To pass the value of `field1` as `test_name` in the YAML settings, you would write:

```yaml
test_name: field1
```

To map the `testList` JSON array to `test_junit_list`, you would specify:

```yaml
test_junit_list: testList
```

For list attributes, the mapping is direct. For example, to map the `class` field in the JSON list to `test_junit_list_class_name`, you would write:

```yaml
test_junit_list_class_name: class
```


### JSON example

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
    image: your_docker_hub_username/harness-junit-converter:latest
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

### Using this JSON example and running directly the plugin in command line

```
go get -u github.com/urfave/cli
go get -u github.com/joho/godotenv
go get -u github.com/pelletier/go-toml
go get -u github.com/sirupsen/logrus
go build -o harness-junit-converter
./harness-junit-converter --json_content='{"test_junit_name": "MyTestSuite", "test_junit_package": "com.example", "test_junit_time": 50, "test_junit_list": [{"test_junit_list_name": "TestCase1", "test_junit_list_class_name": "TestClass1", "test_junit_list_failure": "", "test_junit_list_time": 10}, {"test_junit_list_name": "TestCase2", "test_junit_list_class_name": "TestClass2", "test_junit_list_failure": "Some Failure", "test_junit_list_time": 15}, {"test_junit_list_name": "TestCase3", "test_junit_list_class_name": "TestClass3", "test_junit_list_failure": "", "test_junit_list_time": 25}]}' --test_name="test" --test_description="desc" --test_junit_time="test_junit_time" --test_junit_package="test_junit_package" --test_junit_name="test_junit_name" --test_junit_list="test_junit_list" --test_junit_list_name="test_junit_list_name" --test_junit_list_class_name="test_junit_list_class_name" --test_junit_list_failure="test_junit_list_failure" --test_junit_list_time="test_junit_list_time"



./harness-junit-converter --json_content='[{"code":"DL3018","column":1,"file":"Dockerfile","level":"warning","line":4,"message":"Pin versions in apk add. Instead of `apk add <package>` use `apk add <package>=<version>`"},{"code":"DL3059","column":1,"file":"Dockerfile","level":"info","line":17,"message":"Multiple consecutive `RUN` instructions. Consider consolidation."}]' --test_name="test" --test_description="desc" --test_junit_time="0" --test_junit_package="diego" --test_junit_name="diego-dev" --test_junit_list="." --test_junit_list_name="code" --test_junit_list_class_name="file" --test_junit_list_failure="message" --test_junit_list_time="line"


./harness-junit-converter \
  --json_file_name='tests/nested-json.json' \
  --test_name="object_name" \
  --test_description="file_name" \
  --test_junit_time="file_row" \
  --test_junit_package="file_name" \
  --test_junit_name="object_name" \
  --test_junit_list="checks" \
  --test_junit_list_name="name" \
  --test_junit_list_class_name="comment" \
  --test_junit_list_failure="comments[].summary" \
  --test_junit_list_time="grade" \
  --nested_json_list=true \
  --test_junit_skip_field="skipped
```

## How to Build and Run the Plugin Locally

1. Clone the repository and navigate to its directory.
2. Build the Docker image:

    ```bash
    docker build -t your_docker_hub_username/harness-junit-converter:latest .
    ```

3. Run the Docker container:

    ```bash
    docker run your_docker_hub_username/harness-junit-converter:latest
    ```

4. Push the Docker image to your Docker Hub:

    ```bash
    docker push your_docker_hub_username/harness-junit-converter:latest
    ```

For more details, please refer to the `Dockerfile` in the repository.

## License

This project is licensed under the MIT License.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
