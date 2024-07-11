package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthorizationServiceStackProps struct {
	awscdk.StackProps
}

func NewAuthorizationServiceStack(scope constructs.Construct, id string, props *AuthorizationServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// * lambda handlers
	basicAuthorizer := awslambda.NewFunction(stack, jsii.String("GetProductsHandler"), &awslambda.FunctionProps{
		FunctionName: jsii.String("basicAuthorizer"),
		Code:         awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime:      awslambda.Runtime_NODEJS_18_X(),
		Handler:      jsii.String("basicAuthorizer.handler"),
	})

	// Output the function name to use in other service cdk stack name
	awscdk.NewCfnOutput(stack, jsii.String("function_name"), &awscdk.CfnOutputProps{
		Value: basicAuthorizer.FunctionArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAuthorizationServiceStack(app, "ImportServiceStack", &AuthorizationServiceStackProps{})

	app.Synth(nil)
}
