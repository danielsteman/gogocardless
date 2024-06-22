import getRedirectLink from '../../utils/getRedirectLink';
import { NextResponse } from 'next/server';
import { emailIsValid } from '@/app/utils/emailIsValid';
import authOptions from '@/app/auth';
import { getServerSession } from 'next-auth';
import { redirect } from 'next/navigation';

export async function POST(req: Request) {
  const session = await getServerSession(authOptions);

  if (!session) {
    return NextResponse.redirect("/login")
  }

  const data = await req.json();
  const institutionId = data.institutionId
  const email = session?.user?.email

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
