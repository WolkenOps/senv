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
	client *ssm.SSM
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
	paths := flag.String("paths", "/", "Comma separated paths to look for, i.e. /dev/service, /dev/global/")
	export := flag.Bool("export", false, "The path to look for, i.e. /dev/service")
	flag.Parse()
	parameters, err := fetchParametersByPaths(splitPaths(*paths))
	if err != nil {
		panic(err)
	}
	fmt.Print(formatParameters(parameters, *export))
}

func splitPaths(paths string) []string {
	return strings.Split(paths, ",")
}

func fetchParametersByPaths(paths []string) ([]parameter, error) {
	var parameters []parameter
	for _, path := range paths {
		p, err := fetchParametersByPath(path)
		if err != nil {
			return []parameter{}, err
		}
		parameters = append(parameters, p...)
	}
	return parameters, nil
}

func fetchParametersByPath(path string) ([]parameter, error) {
	var parameters []parameter
	done := false
	var token string
	for !done {
		input := &ssm.GetParametersByPathInput{
			Path:           &path,
			WithDecryption: aws.Bool(true),
		}
		if token != "" {
			input.SetNextToken(token)
		}
		output, err := client.GetParametersByPath(input)
		if err != nil {
			return []parameter{}, err
		}
		for _, p := range output.Parameters {
			name := *p.Name
			value := *p.Value
			if strings.Compare(path, "/") != 0 {
				name = strings.Replace(strings.Trim(name[len(path):], "/"), "/", "_", -1)
			}
			parameters = append(parameters, parameter{name, value})
			if output.NextToken != nil {
				token = *output.NextToken
			} else {
				done = true
			}
		}
	}
	return parameters, nil
}

func formatParameters(parameters []parameter, export bool) string {
	var buffer strings.Builder
	var prefix string
	if export {
		prefix = "export "
	}
	for _, parameter := range parameters {
		buffer.WriteString(fmt.Sprintf("%s%s=%s\n", prefix, parameter.name, parameter.value))
	}
	return buffer.String()
}
