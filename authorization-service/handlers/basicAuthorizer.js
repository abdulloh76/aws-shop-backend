
export async function handler(event) {
  try {
    const authorizationHeader = event.headers.Authorization;

    if (!authorizationHeader) {
      return {
        statusCode: 401,
        body: JSON.stringify({ message: "Unauthorized: Missing Authorization Header" }),
      };
    }

    const encodedCredentials = authorizationHeader.split(' ')[1];
    const decodedCredentials = Buffer.from(encodedCredentials, 'base64').toString('utf-8');
    const [username, password] = decodedCredentials.split('=');

    const expectedPassword = process.env[username];

    if (expectedPassword && expectedPassword === password) {
      return generatePolicy(username, 'Allow', event.methodArn);
    } else {
      return {
        statusCode: 403,
        body: JSON.stringify({ message: "Forbidden: Invalid Credentials" }),
      };
    }
  } catch (error) {
    return {
      statusCode: 500,
      body: JSON.stringify({ message: "Internal Server Error", error: error.message }),
    };
  }
};

const generatePolicy = (principalId, effect, resource) => {
  const policyDocument = {
    Version: '2012-10-17',
    Statement: [{
      Action: 'execute-api:Invoke',
      Effect: effect,
      Resource: resource,
    }],
  };

  return {
    principalId,
    policyDocument,
  };
};
