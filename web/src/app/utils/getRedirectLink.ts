import jwt from 'jsonwebtoken';

export default async function getRedirectLink(
  institutionId: string,
  userEmail: string,
): Promise<Response> {
  const secret = process.env.NEXTAUTH_SECRET;

  if (!secret) {
    throw new Error(
      'NEXTAUTH_SECRET is not defined in the environment variables',
    );
  }

  const payload = {
    email: userEmail,
  };

  const token = jwt.sign(payload, secret, { expiresIn: '1h' });

  try {
    const response = await fetch(
      `http://localhost:3333/api/user/redirect?institutionId=${institutionId}`,
      {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      },
    );
    return response;
  } catch (error) {
    console.error('Error:', error);
    return new Response(
      JSON.stringify({
        message: 'An error occurred while getting the redirect link.',
      }),
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
        },
      },
    );
  }
}
