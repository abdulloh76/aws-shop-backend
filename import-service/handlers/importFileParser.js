import { GetObjectCommand, S3Client, CopyObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";
import csv from 'csv-parser';

const s3Client = new S3Client({ region: process.env.AWS_REGION });

export async function handler(event) {
  try {
    console.log("request:", JSON.stringify(event, undefined, 2));
    const bucketName = process.env.BUCKET_NAME;

    for (const record of event.Records) {
      const objectKey = record.s3.object.key;

      const getCommand = new GetObjectCommand({
        Bucket: bucketName,
        Key: objectKey,
      });

      const { Body } = await s3Client.send(getCommand)
      const products = [];
      Body
        .pipe(csv())
        .on('data', (data) => products.push(data))
        .on('end', () => {
          // * {id,title,description,price,count}
          console.log(products);
          // * just log and thats it)
        });

      const copyCommand = new CopyObjectCommand({
        Bucket: bucketName,
        CopySource: objectKey,
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
