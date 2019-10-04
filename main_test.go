package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockSSMClient struct {
	ssmiface.SSMAPI
}

func (m *mockSSMClient) GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	var output ssm.GetParametersByPathOutput

	if input.Path == aws.String("/") {
		output = ssm.GetParametersByPathOutput{
			NextToken: nil,
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{
					Name:  aws.String("my_parameter"),
					Value: aws.String("my_value"),
				},
				&ssm.Parameter{
					Name:  aws.String("my_other_parameter"),
					Value: aws.String("my_value_value"),
				},
			},
		}
	} else {
		output = ssm.GetParametersByPathOutput{
			NextToken: nil,
			Parameters: []*ssm.Parameter{
				&ssm.Parameter{
					Name:  aws.String("/dev/my_parameter"),
					Value: aws.String("my_value"),
				},
				&ssm.Parameter{
					Name:  aws.String("/dev/my_other_parameter"),
					Value: aws.String("my_value_value"),
				},
			},
		}
	}

	return &output, nil
}

func TestSplitPaths(t *testing.T) {
	total := splitPaths("path1,path2,path2")
	if cap(total) != 3 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", cap(total), 3)
	}
}

func TestFormatParameters(t *testing.T) {
	var value string
	parameterWExport := "export my_parameter=my_value\n"
	parameterWOExport := "my_parameter=my_value\n"
	parameters := map[string]string{
		"my_parameter": "my_value",
	}
	value = formatParameters(parameters, true)
	if parameterWExport != value {
		t.Errorf("Values are different: %s, want: %s.", parameterWExport, value)
	}
	value = formatParameters(parameters, false)
	if parameterWOExport != value {
		t.Errorf("Values are different: %s, want: %s.", parameterWOExport, value)
	}
}

func TestFetchParametersByPaths(t *testing.T) {
	client = &mockSSMClient{}
	paths := "/,/dev/"
        validParameters := []string{
		"/dev/my_parameter",
		"/dev/my_other_parameter",
		"my_parameter",
		"my_other_parameter",
	}

	parameters, err := fetchParametersByPaths(splitPaths(paths))
	if err != nil {
		t.Errorf("an error ocurred %s", err.Error())
	} else {
		for key := range parameters {
			if !contains(validParameters, key) {
				t.Errorf("Contains unexpected values, expected %s", key)
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
