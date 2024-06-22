const { DynamoDBClient, GetItemCommand } = require("@aws-sdk/client-dynamodb");

exports.handler = async function (event) {
  console.log("request:", JSON.stringify(event, undefined, 2));
  const productId = event.pathParameters.productId;
  const client = new DynamoDBClient();

  const getProductCommand = new GetItemCommand({
    TableName: process.env.PRODUCTS_TABLE_NAME,
    Key: { id: { S: productId } },
  });
  const product = await client.send(getProductCommand);
  console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));

  const getStockCommand = new GetItemCommand({
    TableName: process.env.PRODUCTS_TABLE_NAME,
    Key: { product_id: { S: productId } }
  });
  const stock = await client.send(getStockCommand);
  console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));

  if (product) {
    return {
      statusCode: 200,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(product)
    };
  }

  return {
    statusCode: 404,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET",
      "Content-Type": "text/plain",
    },
    body: `Product with id ${productId} not found`
  };
};