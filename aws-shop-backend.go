package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsShopBackendStackProps struct {
	awscdk.StackProps
}

func NewAwsShopBackendStack(scope constructs.Construct, id string, props *AwsShopBackendStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	getProductsHandler := awslambda.NewFunction(stack, jsii.String("GetProductsHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("getProductsList.handler"),
	})

	getProductByIdHandler := awslambda.NewFunction(stack, jsii.String("GetProductByIdHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("getProductsById.handler"),
	})

	productApi := awsapigateway.NewRestApi(stack, jsii.String("Product-Service-Rest-Api"), &awsapigateway.RestApiProps{
		DeployOptions: &awsapigateway.StageOptions{StageName: jsii.String("dev")},
	})

	// /products
	productsResources := productApi.Root().AddResource(jsii.String("products"), &awsapigateway.ResourceOptions{})
	productsResources.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(
			getProductsHandler,
			&awsapigateway.LambdaIntegrationOptions{},
		),
		productsResources.DefaultMethodOptions(),
	)

	// /products/{productId}
	productByIdApiResources := productsResources.AddResource(jsii.String(`{productId}`), &awsapigateway.ResourceOptions{})
	productByIdApiResources.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(
			getProductByIdHandler,
			&awsapigateway.LambdaIntegrationOptions{},
		),
		productByIdApiResources.DefaultMethodOptions(),
	)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAwsShopBackendStack(app, "AwsShopBackendStack", &AwsShopBackendStackProps{})

	app.Synth(nil)
}
