import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { v4 as uuidv4 } from 'uuid';

export async function handler(event) {
  try {
    console.log("🚀 request:", JSON.stringify(event, undefined, 2));
    const { title, description, price, count } = event.body;

    const productId = uuidv4();
    const client = new DynamoDBClient({ region: process.env.AWS_REGION });

    const putProductCommand = new PutItemCommand({
      TableName: process.env.PRODUCTS_TABLE_NAME,
      Item: {
        id: { "S": productId },
        title: { "S": title },
        description: { "S": description },
        price: { "N": price.toString() }
      },
    });
    const product = await client.send(putProductCommand);
    console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));

    const putStockCommand = new PutItemCommand({
      TableName: process.env.STOCKS_TABLE_NAME,
      Item: {
        product_id: { "S": productId },
        count: { "N": count.toString() }
      },
    });
    const stock = await client.send(putStockCommand);
    console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));

    const joinedProduct = {
      id: product.Item.id.S,
      title: product.Item.title.S,
      description: product.Item.description.S,
      price: product.Item.price.N,
      count: stock.Item.count.N
    }

    return {
      statusCode: 200,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(joinedProduct)
    };
  } catch (error) {
    return {
      statusCode: 500,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(error)
    };
  }
}
