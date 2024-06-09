'use client';

import { useSession } from 'next-auth/react';
import { redirect } from 'next/navigation';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

const handleBankClick = async (institutionId: string, userEmail: string) => {
  try {
    const response = await fetch('http://localhost:3333/api/user/redirect', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        institutionId,
        userEmail,
      }),
    });

    if (response.ok) {
      const result = await response.json();
      window.location.href = result.link;
    } else {
      alert('Failed to get redirect link');
    }
  } catch (error) {
    console.error('Error:', error);
    alert('An error occurred while processing your request.');
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
