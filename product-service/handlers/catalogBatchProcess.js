import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { v4 as uuidv4 } from 'uuid';

export async function handler(event) {
  try {
    console.log("🚀 catalogBatchProcess request:", JSON.stringify(event, undefined, 2));
    const consumedProducts = []
    for (const record of event.Records) {
      console.log('SQS Message:', record.body);

      const { title, description, price, count } = JSON.parse(record.body);
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
      consumedProducts.push(joinedProduct);
    }

    console.log("🚀 ~ handler ~ consumedProducts:", consumedProducts);
  } catch (error) {
    console.log("🚀 ~ handler ~ error:", error);
  }
}
