import { GetObjectCommand, S3Client, CopyObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";
import { SQSClient, GetQueueUrlCommand, SendMessageCommand } from "@aws-sdk/client-sqs";
import csv from 'csv-parser';

const s3Client = new S3Client({ region: process.env.AWS_REGION });
const sqsClient = new SQSClient();

export async function handler(event) {
  try {
    console.log("request:", JSON.stringify(event, undefined, 2));
    const bucketName = process.env.BUCKET_NAME;

    const getQueueUrlCommand = new GetQueueUrlCommand({ QueueName: "catalogItemsQueue" });
    const { QueueUrl } = await sqsClient.send(getQueueUrlCommand);

    for (const record of event.Records) {
      const objectKey = record.s3.object.key;

      const getCommand = new GetObjectCommand({
        Bucket: bucketName,
        Key: objectKey,
      });

      const { Body } = await s3Client.send(getCommand)
      const products = await streamParser(Body.pipe(csv()));

      for (const product of products) {
        console.log("🚀 ~ handler ~ product:", product);

        const input = {
          QueueUrl,
          MessageBody: JSON.stringify(product),
        };
        const sendMessageCommand = new SendMessageCommand(input);
        const sqsSendMessageResponse = await sqsClient.send(sendMessageCommand);
        console.log("🚀 ~ handler ~ sqsSendMessageResponse:", sqsSendMessageResponse);
      }

      const copyCommand = new CopyObjectCommand({
        Bucket: bucketName,
        CopySource: `${bucketName}/${objectKey}`,
        Key: `${objectKey.replace('uploaded', 'parsed')}`,
      });
      await s3Client.send(copyCommand)
      const deleteCommand = new DeleteObjectCommand({
        Bucket: bucketName,
        Key: objectKey,
      });
      await s3Client.send(deleteCommand)
    }
  } catch (error) {
    console.log("🚀 ~ handler ~ error:", error);
  }
}

async function streamParser(readableStream) {
  const collectedData = []
  return new Promise((resolve, reject) => {
    readableStream
      .on('data', (data) => collectedData.push(data))
      .on('end', () => {
        // * {id,title,description,price,count}
        resolve(collectedData)
      })
      .on('error', reject);
  });
}
