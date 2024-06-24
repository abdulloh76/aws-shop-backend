import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { v4 as uuidv4 } from 'uuid';
import { generate } from 'random-words';

export async function handler (event) {
  console.log("request:", JSON.stringify(event, undefined, 2));

  const client = new DynamoDBClient({ region: process.env.AWS_REGION });
  const products = [];

  for (let i = 0; i < 30; i++) {
    const productId = uuidv4();
    console.log("🚀 ~ handler ~ productId:", productId);
    const title = generate(1).toString();
    const description = generate(10).join(" ");
    const price = Math.floor(Math.random() * 100 + 1).toString();
    const count = Math.floor(Math.random() * 100 + 1).toString();

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
      TableName: process.env.STOCKS_TABLE_NAME,
      Item: {
        product_id: { "S": productId },
        count: { "N": count }
      },
    });
    const stock = await client.send(putStockCommand);
    console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));

    products.push({ id: productId, title, description, price, count })
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
}
