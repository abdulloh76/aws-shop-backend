import data from '../utils/data.json'

exports.handler = async function(event) {
  console.log("request:", JSON.stringify(event, undefined, 2));
  const productId = event.pathParameters.productId;
  
  const { products } = data;
  const product = products.find(el => el.id == productId)
  
  if (product) {
    return {
      statusCode: 200,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(product)
    };
  }

  return {
    statusCode: 404,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET",
      "Content-Type": "text/plain",
    },
    body: `Product with id ${productId} not found`
  };
};