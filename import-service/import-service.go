package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type ImportServiceStackProps struct {
	awscdk.StackProps
	CatalogQueueArn     string
	CatalogQueueUrl     string
	LambdaAuthorizerArn string
}

func NewImportServiceStack(scope constructs.Construct, id string, props *ImportServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	importsBucket := awss3.NewBucket(stack, jsii.String("products-import"), &awss3.BucketProps{
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{awss3.HttpMethods_GET, awss3.HttpMethods_POST, awss3.HttpMethods_PUT, awss3.HttpMethods_HEAD, awss3.HttpMethods_DELETE},
				AllowedOrigins: jsii.Strings("*"),
				AllowedHeaders: jsii.Strings("*"),
			},
		},
	})

	// * Import the existing SQS Queue by ARN and URL
	catalogQueue := awssqs.Queue_FromQueueAttributes(stack, jsii.String("ImportedQueue"), &awssqs.QueueAttributes{
		QueueUrl: jsii.String(props.CatalogQueueUrl),
		QueueArn: jsii.String(props.CatalogQueueArn),
	})

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
			"BUCKET_NAME":       importsBucket.BucketName(),
			"CATALOG_QUEUE_URL": &props.CatalogQueueUrl,
		},
	})
	lambdaAuthorizerHandler := awslambda.Function_FromFunctionArn(stack, jsii.String("LambdaAuthorizerHandler"), &props.LambdaAuthorizerArn)
	lambdaTokenAuthorizer := awsapigateway.NewTokenAuthorizer(stack, jsii.String("LambdaAuthorizer"), &awsapigateway.TokenAuthorizerProps{
		Handler: lambdaAuthorizerHandler,
	})

	// * queue grant access
	catalogQueue.GrantSendMessages(importFileParserHandler)

	// * bucket grant accesses
	importsBucket.GrantReadWrite(importProductsFileHandler, jsii.String("*"))
	importsBucket.GrantPut(importProductsFileHandler, jsii.String("*"))

	importsBucket.GrantReadWrite(importFileParserHandler, jsii.String("*"))
	importsBucket.GrantPut(importFileParserHandler, jsii.String("*"))
	importsBucket.GrantDelete(importFileParserHandler, jsii.String("*"))

	// * event notifications
	parserNotificationDest := awss3notifications.NewLambdaDestination(importFileParserHandler)
	importsBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, parserNotificationDest, &awss3.NotificationKeyFilter{
		Prefix: jsii.String("uploaded/"),
	})

	// * apigateway instance
	importApi := awsapigateway.NewRestApi(stack, jsii.String("Import-Service-Rest-Api"), &awsapigateway.RestApiProps{
		DeployOptions: &awsapigateway.StageOptions{StageName: jsii.String("dev")},
	})

	importApi.AddGatewayResponse(jsii.String("ResponseType_4XX"), &awsapigateway.GatewayResponseOptions{
		Type: awsapigateway.ResponseType_DEFAULT_4XX(),
		ResponseHeaders: &map[string]*string{
			"Access-Control-Allow-Origin": jsii.String("'*'"),
		},
	})
	importApi.AddGatewayResponse(jsii.String("ResponseType_5XX"), &awsapigateway.GatewayResponseOptions{
		Type: awsapigateway.ResponseType_DEFAULT_5XX(),
		ResponseHeaders: &map[string]*string{
			"Access-Control-Allow-Origin": jsii.String("'*'"),
		},
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
		&awsapigateway.MethodOptions{
			Authorizer: lambdaTokenAuthorizer,
		},
	)

	return stack
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	defer jsii.Close()

	catalogQueueUrl := os.Getenv("CATALOG_QUEUE_URL")
	catalogQueueArn := os.Getenv("CATALOG_QUEUE_ARN")
	lambdaAuthorizerArn := os.Getenv("LAMBDA_AUTHORIZER_NAME")

	app := awscdk.NewApp(nil)

	NewImportServiceStack(app, "ImportServiceStack", &ImportServiceStackProps{
		CatalogQueueArn:     catalogQueueArn,
		CatalogQueueUrl:     catalogQueueUrl,
		LambdaAuthorizerArn: lambdaAuthorizerArn,
	})

	app.Synth(nil)
}
