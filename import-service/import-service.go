package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
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

	importsBucket := awss3.NewBucket(stack, jsii.String("products-import"), &awss3.BucketProps{
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{awss3.HttpMethods_GET, awss3.HttpMethods_POST, awss3.HttpMethods_PUT, awss3.HttpMethods_HEAD, awss3.HttpMethods_DELETE},
				AllowedOrigins: jsii.Strings("*"),
				AllowedHeaders: jsii.Strings("*"),
			},
		},
	})
	importsBucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:     awsiam.Effect_ALLOW,
		Actions:    jsii.Strings("s3:*"),
		Resources:  jsii.Strings(*importsBucket.BucketArn() + "/*"),
		Principals: &[]awsiam.IPrincipal{awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{})},
	}))

	// * lambda handlers
	importProductsFileHandler := awslambda.NewFunction(stack, jsii.String("ImportProductsFileHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("importProductsFile.handler"),
		Environment: &map[string]*string{
			"BUCKET_NAME": importsBucket.BucketName(),
		},
	})
	importFileParserHandler := awslambda.NewFunction(stack, jsii.String("importFileParserHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("importFileParser.handler"),
		Environment: &map[string]*string{
			"BUCKET_NAME": importsBucket.BucketName(),
		},
	})

	// * bucket grant accesses
	importsBucket.GrantReadWrite(importProductsFileHandler, jsii.String("*"))
	importsBucket.GrantPut(importProductsFileHandler, jsii.String("*"))
	importsBucket.GrantReadWrite(importFileParserHandler, jsii.String("*"))
	importsBucket.GrantPut(importFileParserHandler, jsii.String("*"))
	importsBucket.GrantDelete(importFileParserHandler, jsii.String("*"))

	importsBucket.AddToResourcePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:     awsiam.Effect_ALLOW,
		Actions:    jsii.Strings("s3:*"),
		Resources:  jsii.Strings(*importsBucket.BucketArn() + "/*"),
		Principals: &[]awsiam.IPrincipal{importProductsFileHandler.GrantPrincipal(), importFileParserHandler.GrantPrincipal()},
	}))

	// * event notifications
	parserNotificationDest := awss3notifications.NewLambdaDestination(importFileParserHandler)
	importsBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, parserNotificationDest, &awss3.NotificationKeyFilter{
		Prefix: jsii.String("uploaded/"),
	})

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
