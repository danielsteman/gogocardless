import getRedirectLink from '../../utils/getRedirectLink';
import { NextResponse } from 'next/server';
import { emailIsValid } from '@/app/utils/emailIsValid';

export async function POST(req: Request) {
  const data = await req.json();
  console.log(data)
  const email = data.email
  const institutionId = data.institutionId

  if (!institutionId || !email) {
    return NextResponse.json({ error: 'Missing institutionId or userEmail' }, { status: 400 })
  }

  if (!emailIsValid(email)) {
    return NextResponse.json({ error: 'Email has invalid format' }, { status: 400 })
  }

  let response

  try {
    response = await getRedirectLink(institutionId as string, email as string);
    const data = await response.json();
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error('Error:', error);
    return response
  }
}
