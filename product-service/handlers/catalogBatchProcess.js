import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { SNSClient, PublishCommand } from "@aws-sdk/client-sns";
import { v4 as uuidv4 } from 'uuid';

const dynamoDBClient = new DynamoDBClient({ region: process.env.AWS_REGION });
const snsClient = new SNSClient();

export async function handler(event) {
  try {
    console.log("🚀 catalogBatchProcess request:", JSON.stringify(event, undefined, 2));
    const consumedProducts = []

    for (const record of event.Records) {
      const productsFromSQS = JSON.parse(record.body)
      console.log("🚀 ~ handler ~ SQS products:", productsFromSQS);

      for (const productFromMessage of productsFromSQS) {
        const { title, description, price, count } = productFromMessage;
        const productId = uuidv4();

        const putProductCommand = new PutItemCommand({
          TableName: process.env.PRODUCTS_TABLE_NAME,
          Item: {
            id: { "S": productId },
            title: { "S": title },
            description: { "S": description },
            price: { "N": price.toString() }
          },
        });
        const product = await dynamoDBClient.send(putProductCommand);
        console.log("🚀 ~ exports.handler=function ~ product:", JSON.stringify(product));

        const putStockCommand = new PutItemCommand({
          TableName: process.env.STOCKS_TABLE_NAME,
          Item: {
            product_id: { "S": productId },
            count: { "N": count.toString() }
          },
        });
        const stock = await dynamoDBClient.send(putStockCommand);
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
    }

    const emailMessage = {
      Message: 'New Products imported successfully',
      products: consumedProducts
    };
    const publishCommand = new PublishCommand({
      Message: JSON.stringify(emailMessage),
      TopicArn: process.env.CREATE_PRODUCT_TOPIC_ARN,
      MessageAttributes: {
        newProductsAmount: {
          DataType: 'Number',
          StringValue: consumedProducts.length.toString(),
        }
      }
    });
    const snsPublishResponse = await snsClient.send(publishCommand);
    console.log('Message successfully sent to SNS topic: ', snsPublishResponse.MessageId);
  } catch (error) {
    console.log("🚀 ~ handler ~ error:", error);
  }
}
