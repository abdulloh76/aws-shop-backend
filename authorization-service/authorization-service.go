package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/joho/godotenv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthorizationServiceStackProps struct {
	awscdk.StackProps
	SecretKey string
}

func NewAuthorizationServiceStack(scope constructs.Construct, id string, props *AuthorizationServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// * lambda handlers
	basicAuthorizer := awslambda.NewFunction(stack, jsii.String("BasicAuthorizer"), &awslambda.FunctionProps{
		FunctionName: jsii.String("BasicAuthorizer"),
		Code:         awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime:      awslambda.Runtime_NODEJS_18_X(),
		Handler:      jsii.String("basicAuthorizer.handler"),
		Environment: &map[string]*string{
			"SECRET_KEY": &props.SecretKey,
		},
	})

	// Output the function name to use in other service cdk stack name
	awscdk.NewCfnOutput(stack, jsii.String("function_name"), &awscdk.CfnOutputProps{
		Value: basicAuthorizer.FunctionArn(),
	})

	return stack
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	defer jsii.Close()

	secretKey := os.Getenv("SECRET_KEY")

	app := awscdk.NewApp(nil)

	NewAuthorizationServiceStack(app, "ImportServiceStack", &AuthorizationServiceStackProps{
		SecretKey: secretKey,
	})

	app.Synth(nil)
}
