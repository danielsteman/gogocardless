'use client';

import { redirect } from 'next/navigation';
import { useSession } from 'next-auth/react';
import jwt from 'jsonwebtoken';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

const handleBankClick = async (institutionId: string, email: string) => {
  try {
    var token = jwt.sign({ email }, process.env.NEXT_PUBLIC_NEXTAUTH_SECRET!);
    const response = await fetch('/api/redirect', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: token,
      },
      body: JSON.stringify({
        institutionId: institutionId,
        email: email,
      }),
    });

    const data = await response.json();
    if (response.ok) {
      localStorage.setItem('agreementRef', data.id);
      window.location.href = data.link;
    } else {
      console.error('Error:', data.error);
    }
  } catch (error) {
    console.error('Error:', error);
  }
};

export default function BanksList({ banks }: { banks: Bank[] }) {
  const { data: session, status } = useSession();

  if (status === 'loading') {
    return <p>Loading...</p>;
  }

  if (!session || !session.user) {
    return <p>You need to be authenticated to view this content.</p>;
  }

  if (!session.user.email) {
    console.warn('could not find email in session, redirecting to login page');
    redirect('/login');
  }

  const userEmail = session.user.email;

  return (
    <div>
      <h1 className='text-2xl pb-8 font-bold'>Banks</h1>
      <div className='flex flex-col items-start'>
        {banks.map(bank => (
          <button
            key={bank.id}
            onClick={() => handleBankClick(bank.id, userEmail)}
            className='text-left mb-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700'
          >
            {bank.name}
          </button>
        ))}
      </div>
    </div>
  );
}
