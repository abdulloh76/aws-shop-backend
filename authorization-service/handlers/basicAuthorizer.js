
export async function handler(event) {
  console.log("🚀 ~ handler ~ event:", event);
  try {
    const authorizationHeader = event.authorizationToken;
    console.log("🚀 ~ handler ~ authorizationHeader:", authorizationHeader);

    if (!authorizationHeader) {
      return {
        statusCode: 401,
        body: JSON.stringify({ message: "Unauthorized: Missing Authorization Header" }),
      };
    }

    const encodedCredentials = authorizationHeader.split(' ')[1];
    const decodedCredentials = Buffer.from(encodedCredentials, 'base64').toString('utf-8');
    const [username, password] = decodedCredentials.split('=');
    console.log("🚀 ~ handler ~ username:", username);
    console.log("🚀 ~ handler ~ password:", password);

    const [_, expectedPassword] = process.env.SECRET_KEY?.split('=')
    console.log("🚀 ~ handler ~ expectedPassword:", expectedPassword);

    if (expectedPassword && expectedPassword === password) {
      return generatePolicy(username, 'Allow', event.methodArn);
    } else {
      return generatePolicy(username, 'Deny', event.methodArn);
    }
  } catch (error) {
    console.log("🚀 ~ handler ~ error:", error);
    return {
      statusCode: 500,
      body: JSON.stringify({ message: "Internal Server Error", error: error.message }),
    };
  }
};

const generatePolicy = (principalId, effect, resource,) => {
  return {
    principalId,
    policyDocument: {
      Version: '2012-10-17',
      Statement: [
        {
          Action: ['execute-api:Invoke'],
          Effect: effect,
          Resource: resource
        }
      ]
    }
  }
};
