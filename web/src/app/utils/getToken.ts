import jwt from 'jsonwebtoken';

export async function getToken(email: string): Promise<string> {
  const token = jwt.sign({ email }, process.env.NEXTAUTH_SECRET!, {
    expiresIn: '1h',
  });
  return token
}