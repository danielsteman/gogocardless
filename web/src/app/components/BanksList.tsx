'use client';

import { useState } from 'react';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function handleBankClick(id: string) {
  const response = await fetch(
    `http://localhost:3333/api/user/redirect?institutionId=${id}`,
    {
      method: 'GET',
    },
  );

  if (response.ok) {
    const result = await response.json();
    alert(result.message);
  } else {
    alert('Failed to handle click');
  }
}

export default function BanksList({ banks }: { banks: Bank[] }) {
  return (
    <div className='px-16 py-8'>
      <h1 className='text-2xl pb-8 font-bold'>Banks</h1>
      <div className='flex flex-col items-start'>
        {banks.map(bank => (
          <button
            key={bank.id}
            onClick={() => handleBankClick(bank.id)}
            className='text-left mb-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700'
          >
            {bank.name}
          </button>
        ))}
      </div>
    </div>
  );
}
