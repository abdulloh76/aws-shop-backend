package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type ProductServiceStackProps struct {
	awscdk.StackProps
	SubscriptionEmailAddress string
}

func NewProductServiceStack(scope constructs.Construct, id string, props *ProductServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// * dynamodb tables
	productsTable := awsdynamodb.NewTable(stack, jsii.String("products"), &awsdynamodb.TableProps{
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("id"), Type: awsdynamodb.AttributeType_STRING},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	stocksTable := awsdynamodb.NewTable(stack, jsii.String("stocks"), &awsdynamodb.TableProps{
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("product_id"), Type: awsdynamodb.AttributeType_STRING},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// * amazon sqs
	catalogItemsQueue := awssqs.NewQueue(stack, jsii.String("catalogItemsQueue"), &awssqs.QueueProps{
		QueueName: jsii.String("catalogItemsQueue"),
	})
	// * aws sns - new topic and email subscription
	createProductTopic := awssns.NewTopic(stack, jsii.String("createProductTopic"), &awssns.TopicProps{})
	createProductTopic.AddSubscription(awssnssubscriptions.NewEmailSubscription(&props.SubscriptionEmailAddress, &awssnssubscriptions.EmailSubscriptionProps{
		FilterPolicy: &map[string]awssns.SubscriptionFilter{
			"newProductsAmount": awssns.SubscriptionFilter_NumericFilter(&awssns.NumericConditions{
				GreaterThan: jsii.Number(3),
			}),
		},
	}))

	// * lambda handlers
	getProductsHandler := awslambda.NewFunction(stack, jsii.String("GetProductsHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("getProductsList.handler"),
		Environment: &map[string]*string{
			"PRODUCTS_TABLE_NAME": productsTable.TableName(),
			"STOCKS_TABLE_NAME":   stocksTable.TableName(),
		},
	})
	getProductByIdHandler := awslambda.NewFunction(stack, jsii.String("GetProductByIdHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("getProductsById.handler"),
		Environment: &map[string]*string{
			"PRODUCTS_TABLE_NAME": productsTable.TableName(),
			"STOCKS_TABLE_NAME":   stocksTable.TableName(),
		},
	})
	createProductHandler := awslambda.NewFunction(stack, jsii.String("CreateProductHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("createProduct.handler"),
		Environment: &map[string]*string{
			"PRODUCTS_TABLE_NAME": productsTable.TableName(),
			"STOCKS_TABLE_NAME":   stocksTable.TableName(),
		},
	})
	generateRandomStockProductsLambda := awslambda.NewFunction(stack, jsii.String("GenerateRandomStockProductsLambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("migrateRandomData.handler"),
		Environment: &map[string]*string{
			"PRODUCTS_TABLE_NAME": productsTable.TableName(),
			"STOCKS_TABLE_NAME":   stocksTable.TableName(),
		},
	})
	catalogBatchProcessHandler := awslambda.NewFunction(stack, jsii.String("CatalogBatchProcessHandler"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("handlers"), nil),
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("catalogBatchProcess.handler"),
		Environment: &map[string]*string{
			"PRODUCTS_TABLE_NAME":      productsTable.TableName(),
			"STOCKS_TABLE_NAME":        stocksTable.TableName(),
			"CREATE_PRODUCT_TOPIC_ARN": createProductTopic.TopicArn(),
		},
	})

	// * lambda event source
	awslambda.NewEventSourceMapping(stack, jsii.String("sqsTrigger"), &awslambda.EventSourceMappingProps{
		Target:         catalogBatchProcessHandler,
		EventSourceArn: catalogItemsQueue.QueueArn(),
		BatchSize:      jsii.Number(5),
	})

	// * queue grant access
	catalogItemsQueue.GrantConsumeMessages(catalogBatchProcessHandler)
	createProductTopic.GrantPublish(catalogBatchProcessHandler)

	// * dynamodb table grant access
	productsTable.GrantReadWriteData(getProductsHandler)
	productsTable.GrantReadWriteData(getProductByIdHandler)
	productsTable.GrantReadWriteData(createProductHandler)
	productsTable.GrantReadWriteData(generateRandomStockProductsLambda)
	productsTable.GrantReadWriteData(catalogBatchProcessHandler)
	stocksTable.GrantReadWriteData(getProductsHandler)
	stocksTable.GrantReadWriteData(getProductByIdHandler)
	stocksTable.GrantReadWriteData(createProductHandler)
	stocksTable.GrantReadWriteData(generateRandomStockProductsLambda)
	stocksTable.GrantReadWriteData(catalogBatchProcessHandler)

	// * apigateway instance
	productApi := awsapigateway.NewRestApi(stack, jsii.String("Product-Service-Rest-Api"), &awsapigateway.RestApiProps{
		DeployOptions: &awsapigateway.StageOptions{StageName: jsii.String("dev")},
	})

	// * /products - GET
	productsResources := productApi.Root().AddResource(jsii.String("products"), &awsapigateway.ResourceOptions{})
	productsResources.AddMethod(
		jsii.String("GET"),
		awsapigateway.NewLambdaIntegration(
			getProductsHandler,
			&awsapigateway.LambdaIntegrationOptions{},
		),
		productsResources.DefaultMethodOptions(),
	)

	// * /products - POST
	productsResources.AddMethod(
		jsii.String("POST"),
		awsapigateway.NewLambdaIntegration(
			createProductHandler,
			&awsapigateway.LambdaIntegrationOptions{},
		),
		productsResources.DefaultMethodOptions(),
	)

	// * /products/{productId}
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
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	defer jsii.Close()

	emailAddress := os.Getenv("EMAIL_ADDRESS")

	app := awscdk.NewApp(nil)

	NewProductServiceStack(app, "ProductServiceStack", &ProductServiceStackProps{
		SubscriptionEmailAddress: emailAddress,
	})

	app.Synth(nil)
}
