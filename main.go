package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	client     *ssm.SSM
	parameters []parameter
)

func init() {
	session := session.Must(session.NewSession())
	client = ssm.New(session)
}

type parameter struct {
	name  string
	value string
}

func main() {
	path := flag.String("path", "/", "The path to look for, i.e. /dev/service")
	export := flag.Bool("export", false, "The path to look for, i.e. /dev/service")
	flag.Parse()

	parameters, err := fetchParameters(*path, "")
	if err != nil {
		panic(err)
	}

	fmt.Print(formatParameters(parameters, *export))
}

func fetchParameters(path string, token string) ([]parameter, error) {
	var parameters []parameter

	input := &ssm.GetParametersByPathInput{
		Path:           &path,
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(false),
	}

	if token != "" {
		input.SetNextToken(token)
	}

	output, err := client.GetParametersByPath(input)

	if err != nil {
		return nil, err
	}

	for _, p := range output.Parameters {
		name := *p.Name
		value := *p.Value

		if strings.Compare(path, "/") != 0 {
			name = strings.Replace(strings.Trim(name[len(path):], "/"), "/", "_", -1)
		}
		parameters = append(parameters, parameter{name, value})
	}

	if output.NextToken != nil {
		parameters, err = fetchParameters(path, *output.NextToken)
	}
	return parameters, nil
}

func formatParameters(parameters []parameter, export bool) string {
	var buffer strings.Builder
	var processed string
	for _, parameter := range parameters {
		if export {
			processed = fmt.Sprintf("export %s=%s\n", parameter.name, parameter.value)
		} else {
			processed = fmt.Sprintf("%s=%s\n", parameter.name, parameter.value)
		}
		buffer.WriteString(processed)
	}
	return buffer.String()
}
