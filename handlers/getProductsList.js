import data from '../utils/data.json'

exports.handler = async function(event) {
  const { products } = data

  console.log("request:", JSON.stringify(event, undefined, 2));

  return {
    statusCode: 200,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(products)
  };
};