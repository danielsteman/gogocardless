import { getToken } from '../utils/getToken';
import { getServerSession } from 'next-auth';
import authOptions from '../auth';
import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import AccountsList from './AccountsList';

async function fetchAccounts(email: string, ref: string): Promise<string[]> {
  const token = await getToken(email);
  const response = await fetch(
    `http://localhost:3333/api/user/accounts?agreementRef=${ref}`,
    {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    },
  );
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(
      `Error fetching accounts: ${response.status} - ${errorText}`,
    );
  }
  const data = await response.json();
  return data.accounts;
}

export default async function Page() {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect('/login');
  }

  const refCookie = cookies().get('agreementRef');
  if (!refCookie) {
    console.error('Agreement reference not found in URL');
  }

  const ref = refCookie!.value;

  let accounts: string[] = [];

  try {
    accounts = await fetchAccounts(session.user?.email!, ref);
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  return <AccountsList accounts={accounts} />;
}
