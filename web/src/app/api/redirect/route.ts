import type { NextApiRequest, NextApiResponse } from 'next';
import getRedirectLink from '../../utils/getRedirectLink';

export async function GET(req: NextApiRequest, res: NextApiResponse) {
  const { institutionId, userEmail } = req.query;

  if (!institutionId || !userEmail) {
    return res.status(400).json({ error: 'Missing institutionId or userEmail' });
  }

  try {
    const response = await getRedirectLink(institutionId as string, userEmail as string);
    const data = await response.json();
    return res.status(response.status).json(data);
  } catch (error) {
    console.error('Error:', error);
    return res.status(500).json({ error: 'An error occurred while getting the redirect link' });
  }
}
