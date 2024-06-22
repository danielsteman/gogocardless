'use client';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

const handleBankClick = async (institutionId: string) => {
  try {
    const response = await fetch('/api/redirect', {
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
      localStorage.setItem('agreementRef', data.id);
      window.location.href = data.link;
    } else {
      console.error('Response error:', data.error);
    }
  } catch (error) {
    console.error('Error:', error);
  }
};

export default function BanksList({ banks }: { banks: Bank[] }) {
  return (
    <div>
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
