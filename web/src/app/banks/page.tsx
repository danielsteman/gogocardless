import { Suspense } from 'react';
import BanksList from './BanksList';
import { redirect } from 'next/navigation';
import { getServerSession } from 'next-auth/next';
import authOptions from '../auth';
import ProtectedLayout from '../layouts/ProtectedLayout';
import jwt from 'jsonwebtoken';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function fetchBanks(email: string): Promise<Bank[]> {
  const token = jwt.sign({ email: email }, process.env.NEXTAUTH_SECRET!, {
    expiresIn: '1h',
  });
  const response = await fetch('http://localhost:3333/api/banks/list', {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Error fetching banks: ${response.status} - ${errorText}`);
  }
  return response.json();
}

export default async function Page() {
  const session = await getServerSession(authOptions);

  if (!session) {
    redirect('/login');
  }

  let banks: Bank[] = [];

  try {
    banks = await fetchBanks(session.user?.email!);
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <ProtectedLayout>
        <p>Welcome, {session?.user?.name}</p>
        <BanksList banks={banks} />
      </ProtectedLayout>
    </Suspense>
  );
}
