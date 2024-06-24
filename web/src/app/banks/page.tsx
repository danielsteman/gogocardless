import { Suspense } from 'react';
import { redirect } from 'next/navigation';
import { getServerSession } from 'next-auth/next';
import authOptions from '../auth';
import ProtectedLayout from '../layouts/ProtectedLayout';
import jwt from 'jsonwebtoken';
import { cookies } from 'next/headers';
import { getToken } from '../utils/getToken';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function fetchBanks(email: string): Promise<Bank[]> {
  const token = await getToken(email);
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

interface Requisition {
  id: string;
  redirect: string;
  status: string;
  agreement: string;
  accounts: string[];
  reference: string;
  userLanguage: string;
  link: string;
}

async function fetchRedirectLink(
  email: string,
  institutionId: string,
): Promise<Requisition> {
  const token = await getToken(email);
  const response = await fetch(
    `http://localhost:3333/api/user/redirect?institutionId=${institutionId}`,
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
      `Error fetching redirect link: ${response.status} - ${errorText}`,
    );
  }
  return response.json();
}

export default async function Page() {
  const session = await getServerSession(authOptions);

  if (!session) {
    redirect('/login');
  }

  const email = session.user?.email!;
  let banks: Bank[] = [];

  try {
    banks = await fetchBanks(email);
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  const handleBankClick = async (data: FormData) => {
    'use server';

    const institutionId = data.get('institutionId') as string;
    let redirectPath: string | null = null;

    try {
      const data = await fetchRedirectLink(email, institutionId);
      if (data) {
        cookies().set('agreementRef', data.id);
        redirectPath = data.link;
      } else {
        redirectPath = '/';
        console.error('Failed to get redirect link');
      }
    } catch (error) {
      redirectPath = '/';
      console.error('Error while getting redirect link:', error);
    } finally {
      if (redirectPath) {
        redirect(redirectPath);
      }
    }
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <ProtectedLayout>
        <p>Welcome, {session?.user?.name}</p>
        <div>
          <h1 className='text-2xl pb-8 font-bold'>Banks</h1>
          <div className='flex flex-col items-start'>
            {banks.map(bank => (
              <form action={handleBankClick} key={bank.id}>
                <input
                  name='institutionId'
                  className='hidden'
                  value={bank.id}
                  readOnly
                />
                <button
                  className='text-left mb-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700'
                  type='submit'
                >
                  {bank.name}
                </button>
              </form>
            ))}
          </div>
        </div>
      </ProtectedLayout>
    </Suspense>
  );
}
