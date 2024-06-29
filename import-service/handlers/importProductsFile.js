import { PutObjectCommand, S3Client } from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";

export async function handler(event) {
  try {
    console.log("request:", JSON.stringify(event, undefined, 2));
    const fileName = event.queryParameters.name;
    console.log("🚀 ~ handler ~ fileName:", fileName);

    const s3Client = new S3Client({ region: process.env.AWS_REGION });
    const command = new PutObjectCommand({
      Bucket: process.env.BUCKET_NAME,
      Key: fileName,
    });
    
    const presignedUrl = getSignedUrl(s3Client, command);
    console.log("🚀 ~ handler ~ presignedUrl:", presignedUrl);

    return {
      statusCode: 200,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: presignedUrl
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