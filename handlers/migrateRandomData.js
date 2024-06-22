const { DynamoDBClient, PutItemCommand } = require("@aws-sdk/client-dynamodb");
const { v4: uuidv4 } = require('uuid');
const { generate } = require('random-words');

exports.handler = async function (event) {
  console.log("request:", JSON.stringify(event, undefined, 2));

  const client = new DynamoDBClient();
  const products = [];

  for (let i = 0; i < 30; i++) {
    const productId = uuidv4();

    const putProductCommand = new PutItemCommand({
      TableName: process.env.PRODUCTS_TABLE_NAME,
      Item: {
        id: { "S": productId },
        title: { "S": generate(1) },
        description: { "S": generate(10).join(" ") },
        price: { "N": Math.floor(Math.random() * 100 + 1) }
      },
    });
    const product = await client.send(putProductCommand);
    console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));
    product.push(product)

    const putStockCommand = new PutItemCommand({
      TableName: process.env.PRODUCTS_TABLE_NAME,
      Item: {
        product_id: { "S": productId },
        count: { "N": Math.floor(Math.random() * 100 + 1) }
      },
    });
    const stock = await client.send(putStockCommand);
    console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));
  }

  return {
    statusCode: 200,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(products)
  };
};