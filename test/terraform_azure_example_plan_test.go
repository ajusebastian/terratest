package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// An example of how to test the Terraform module in examples/terraform-aws-example using Terratest.
func TestTerraformAzureExamplePlan(t *testing.T) {
	t.Parallel()

	// Make a copy of the terraform module to a temporary directory. This allows running multiple tests in parallel
	// against the same terraform module.
	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-azure-example")

	// Give this EC2 Instance a unique ID for a name tag so we can distinguish it from any other EC2 Instance running
	// in your AWS account
	expectedName := fmt.Sprintf("%s", "storageaccountaju1234")
	//expectedName1 := "ajustorageAccount1234123"
	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	//awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	//instanceType := aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro", "t3.micro"})

	// website::tag::1::Configure Terraform setting path to Terraform code, EC2 instance name, and AWS Region. We also
	// configure the options with default retryable errors to handle the most common retryable errors encountered in
	// terraform testing.
	planFilePath := filepath.Join(exampleFolder, "plan.out")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-azure-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"storageaccountname": expectedName,
			//"instance_type": instanceType,
		},

		// Environment variables to set when running Terraform
		/*EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},*/

		// Configure a plan file path so we can introspect the plan and make assertions about it.
		PlanFilePath: planFilePath,
	})

	// website::tag::2::Run `terraform init`, `terraform plan`, and `terraform show` and fail the test if there are any errors
	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	// website::tag::3::Use the go struct to introspect the plan values.
	terraform.RequirePlannedValuesMapKeyExists(t, plan, "azurerm_storage_account.aju-storageaccount")
	azureResource := plan.ResourcePlannedValuesMap["azurerm_storage_account.aju-storageaccount"]
	azurestoreagename := azureResource.AttributeValues["name"]
	assert.Equal(t, expectedName, azurestoreagename)
	//azuretags := azureResource.AttributeValues["tags"].(map[string]interface{})
	//assert.Equal(t, map[string]interface{}{"Name": expectedName}, azuretags)

	// website::tag::4::Alternatively, you can get the direct JSON output and use jsonpath to extract the data.
	// jsonpath only returns lists.
	/*var jsonEC2Tags []map[string]interface{}
	jsonOut := terraform.InitAndPlanAndShow(t, terraformOptions)
	k8s.UnmarshalJSONPath(
		t,
		[]byte(jsonOut),
		"{ .planned_values.root_module.resources[0].values.tags }",
		&jsonEC2Tags,
	)
	assert.Equal(t, map[string]interface{}{"Name": expectedName}, jsonEC2Tags[0])*/
}
