import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { v4 as uuidv4 } from 'uuid';

export async function handler(event) {
  try {
    console.log("🚀 request:", JSON.stringify(event, undefined, 2));
    const { title, description, price, count } = JSON.parse(event.body);

    if (!title || !description || !price || !count) {
      return {
        statusCode: 400,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Methods": "GET",
          "Content-Type": "application/json",
        },
        body: "one or more fields are missing for product creation"
      };
    }

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
      id: productId,
      title: title,
      description: description,
      price: price,
      count: count
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
    console.log("🚀 ~ handler ~ error:", error);
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
