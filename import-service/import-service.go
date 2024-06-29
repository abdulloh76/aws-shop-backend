package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ImportServiceStackProps struct {
	awscdk.StackProps
}

func NewImportServiceStack(scope constructs.Construct, id string, props *ImportServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	shopImports := awss3.NewBucket(stack, jsii.String("products-import"), &awss3.BucketProps{})

	// * lambda handlers
	importProductsFileHandler := awslambda.NewFunction(stack, jsii.String("ImportProductsFileHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("importProductsFile.handler"),
		Environment: &map[string]*string{
			"BUCKET_NAME": shopImports.BucketName(),
		},
	})

	shopImports.GrantReadWrite(importProductsFileHandler, importProductsFileHandler.Role())

	// * apigateway instance
	importApi := awsapigateway.NewRestApi(stack, jsii.String("Import-Service-Rest-Api"), &awsapigateway.RestApiProps{
		DeployOptions: &awsapigateway.StageOptions{StageName: jsii.String("dev")},
	})

	// The name will be passed in a query string as a name parameter and should be described in the AWS CDK Stack as a request parameter.
	// * /import - GET
	importsResources := importApi.Root().AddResource(jsii.String("import"), &awsapigateway.ResourceOptions{})
	importsResources.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(
			importProductsFileHandler,
			&awsapigateway.LambdaIntegrationOptions{},
		),
		importsResources.DefaultMethodOptions(),
	)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewImportServiceStack(app, "ImportServiceStack", &ImportServiceStackProps{})

	app.Synth(nil)
}