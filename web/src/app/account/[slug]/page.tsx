import authOptions from '@/app/auth';
import { getToken } from '@/app/utils/getToken';
import { getServerSession } from 'next-auth';
import { redirect } from 'next/navigation';

async function fetchTransactions(email: string, accountId: string) {
  const token = await getToken(email);
  const response = await fetch(
    `http://localhost:3333/api/user/accounts/${accountId}`,
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
      `Error fetching transactions: ${response.status} - ${errorText}`,
    );
  }
  const data = await response.json();
  return data;
}

export default async function Page({ params }: { params: { slug: string } }) {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect('/login');
  }

  const accountId = params.slug;
  const transactions = await fetchTransactions(session.user?.email!, accountId);

  return <div>Account: {params.slug}</div>;
}
