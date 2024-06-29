import { DynamoDBClient, GetItemCommand } from "@aws-sdk/client-dynamodb";

export async function handler(event) {
  try {
    console.log("request:", JSON.stringify(event, undefined, 2));
    const productId = event.pathParameters.productId;
    const client = new DynamoDBClient({ region: process.env.AWS_REGION });

    const getProductCommand = new GetItemCommand({
      TableName: process.env.PRODUCTS_TABLE_NAME,
      Key: { id: { S: productId } },
    });
    const product = await client.send(getProductCommand);
    console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));

    if (!product) {
      return {
        statusCode: 404,
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Methods": "GET",
          "Content-Type": "text/plain",
        },
        body: `Product with id ${productId} not found`
      };
    }

    const getStockCommand = new GetItemCommand({
      TableName: process.env.STOCKS_TABLE_NAME,
      Key: { product_id: { S: productId } }
    });
    const stock = await client.send(getStockCommand);
    console.log("🚀 ~ exports.handler=function ~ stock:", JSON.stringify(stock));

    const joinedProduct = {
      id: product.Item.id.S,
      title: product.Item.title.S,
      description: product.Item.description.S,
      price: Number(product.Item.price.N),
      count: Number(stock.Item.count.N)
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
