import { DynamoDBClient, ScanCommand } from "@aws-sdk/client-dynamodb";

export async function handler(event) {
  try {
    console.log("request:", JSON.stringify(event, undefined, 2));
    const client = new DynamoDBClient({ region: process.env.AWS_REGION });

    const getProductsCommand = new ScanCommand({
      TableName: process.env.PRODUCTS_TABLE_NAME
    });
    const products = await client.send(getProductsCommand);
    console.log("🚀 ~ exports.handler=function ~ products:", JSON.stringify(products));

    const getStocksCommand = new ScanCommand({
      TableName: process.env.STOCKS_TABLE_NAME,
    });
    const stocks = await client.send(getStocksCommand);
    console.log("🚀 ~ exports.handler=function ~ stocks:", JSON.stringify(stocks));

    const joinedProducts = products.Items.map((p) => {
      const s = stocks.Items.find(s => s.product_id.S === p.id.S)
      return ({
        id: p.id.S,
        title: p.title.S,
        description: p.description.S,
        price: p.price.N,
        count: s.count.N
      })
    })

    return {
      statusCode: 200,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(joinedProducts)
    };
  } catch (error) {
    return {
      statusCode: 500,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: error
    };
  }
}
