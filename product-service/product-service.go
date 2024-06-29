package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ProductServiceStackProps struct {
	awscdk.StackProps
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

	productsTable.GrantReadWriteData(getProductsHandler)
	productsTable.GrantReadWriteData(getProductByIdHandler)
	productsTable.GrantReadWriteData(createProductHandler)
	productsTable.GrantReadWriteData(generateRandomStockProductsLambda)
	stocksTable.GrantReadWriteData(getProductsHandler)
	stocksTable.GrantReadWriteData(getProductByIdHandler)
	stocksTable.GrantReadWriteData(createProductHandler)
	stocksTable.GrantReadWriteData(generateRandomStockProductsLambda)

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
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewProductServiceStack(app, "ProductServiceStack", &ProductServiceStackProps{})

	app.Synth(nil)
}
