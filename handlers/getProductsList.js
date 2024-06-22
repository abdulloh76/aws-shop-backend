const { DynamoDBClient, BatchGetItemCommand } = require("@aws-sdk/client-dynamodb");

exports.handler = async function (event) {
  console.log("request:", JSON.stringify(event, undefined, 2));
  const client = new DynamoDBClient();

  const getProductsCommand = new BatchGetItemCommand({
    RequestItems: process.env.PRODUCTS_TABLE_NAME,
  });
  const products = await client.send(getProductsCommand);
  console.log("🚀 ~ exports.handler=function ~ products:", JSON.stringify(products));

  const getStocksCommand = new BatchGetItemCommand({
    RequestItems: process.env.STOCKS_TABLE_NAME,
  });
  const stocks = await client.send(getStocksCommand);
  console.log("🚀 ~ exports.handler=function ~ stocks:", JSON.stringify(stocks));

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