package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/lambda"

	res "github.com/cloudetc/awsweeper/resource"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcc_LambdaFunction_DeleteByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test.")
	}

	env := InitEnv(t)

	terraformDir := "./test-fixtures/lambda-function"

	terraformOptions := getTerraformOptions(terraformDir, env)

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	id := terraform.Output(t, terraformOptions, "id")

	assertLambdaFunctionExists(t, id)

	writeConfigID(t, terraformDir, res.LambdaFunction, id)
	defer os.Remove(terraformDir + "/config.yml")

	logBuffer, err := runBinary(t, terraformDir, "YES\n")
	require.NoError(t, err)

	assertLambdaFunctionDeleted(t, id)

	fmt.Println(logBuffer)
}

func TestAcc_LambdaFunction_DeleteByTag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test.")
	}

	env := InitEnv(t)

	terraformDir := "./test-fixtures/lambda-function"

	terraformOptions := getTerraformOptions(terraformDir, env)

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	id := terraform.Output(t, terraformOptions, "id")

	assertLambdaFunctionExists(t, id)

	writeConfigTag(t, terraformDir, res.LambdaFunction)
	defer os.Remove(terraformDir + "/config.yml")

	logBuffer, err := runBinary(t, terraformDir, "YES\n")
	require.NoError(t, err)

	assertLambdaFunctionDeleted(t, id)

	fmt.Println(logBuffer)
}

func assertLambdaFunctionExists(t *testing.T, id string) {
	assert.True(t, lambdaFunctionExists(t, id))
}

func assertLambdaFunctionDeleted(t *testing.T, id string) {
	assert.False(t, lambdaFunctionExists(t, id))
}

func lambdaFunctionExists(t *testing.T, id string) bool {
	conn := sharedAwsClient.LambdaAPI

	opts := &lambda.GetFunctionInput{
		FunctionName: &id,
	}

	_, err := conn.GetFunction(opts)
	if err != nil {
		elbErr, ok := err.(awserr.Error)
		if !ok {
			t.Fatal(err)
		}
		if elbErr.Code() == "ResourceNotFoundException" {
			return false
		}
		t.Fatal(err)
	}

	return true
}
