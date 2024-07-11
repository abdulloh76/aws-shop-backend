package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
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

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAuthorizationServiceStack(app, "ImportServiceStack", &AuthorizationServiceStackProps{})

	app.Synth(nil)
}
