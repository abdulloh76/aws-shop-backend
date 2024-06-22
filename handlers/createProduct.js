const { DynamoDBClient, PutItemCommand } = require("@aws-sdk/client-dynamodb");
const { v4: uuidv4 } = require('uuid');

exports.handler = async function (event) {
  console.log("request:", JSON.stringify(event, undefined, 2));
  const { title, description, price, count } = event.params.body;

  const productId = uuidv4();
  const client = new DynamoDBClient();

  const putProductCommand = new PutItemCommand({
    TableName: process.env.PRODUCTS_TABLE_NAME,
    Item: {
      id: { "S": productId },
      title: { "S": title },
      description: { "S": description },
      price: { "N": price }
    },
  });
  const product = await client.send(putProductCommand);
  console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));

  const putStockCommand = new PutItemCommand({
    TableName: process.env.PRODUCTS_TABLE_NAME,
    Item: {
      product_id: { "S": productId },
      count: { "N": count }
    },
  });
  const stock = await client.send(putStockCommand);
  console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));

  return {
    statusCode: 200,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(product)
  };
};