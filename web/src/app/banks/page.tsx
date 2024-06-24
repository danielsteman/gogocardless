import { Suspense } from 'react';
import { redirect } from 'next/navigation';
import { getServerSession } from 'next-auth/next';
import authOptions from '../auth';
import ProtectedLayout from '../layouts/ProtectedLayout';
import jwt from 'jsonwebtoken';
import { cookies } from 'next/headers';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function fetchBanks(email: string): Promise<Bank[]> {
  const token = jwt.sign({ email }, process.env.NEXTAUTH_SECRET!, {
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

  const handleBankClick = async (data: FormData) => {
    'use server';

    console.log(data);
    const institutionId = data.get('institutionId');

    try {
      const response = await fetch('http://localhost:3000/api/redirect', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          institutionId: institutionId,
        }),
      });

      const data = await response.json();
      if (response.ok) {
        cookies().set('agreementRef', data.id);
        window.location.href = data.link;
      } else {
        console.error('Response error:', data.error);
      }
    } catch (error) {
      console.error('Error:', error);
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
