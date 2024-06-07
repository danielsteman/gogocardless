import { Suspense } from 'react';
import BanksList from '../components/BanksList';
import { useSession } from 'next-auth/react';
import { redirect } from 'next/navigation';
import { getServerSession } from 'next-auth/next';
import authOptions from '../auth';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function fetchBanks(): Promise<Bank[]> {
  const response = await fetch('http://localhost:3333/api/banks/list');
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
}

export default async function Page() {
  const session = await getServerSession(authOptions);

  if (!session) {
    redirect('/login');
    return null;
  }

  let banks: Bank[] = [];

  try {
    banks = await fetchBanks();
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <p>Welcome, {session?.user?.name}</p>
        <BanksList banks={banks} />
      </Suspense>
    </>
  );
}
