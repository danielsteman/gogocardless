import { Suspense } from 'react';

interface Bank {
  id: string;
  name: string;
  bic: string;
  transaction_total_days: string;
  countries: string[];
  logo: string;
}

async function fetchBanks(): Promise<Bank[]> {
  const response = await fetch('http://localhost:3333/banks/list');
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  return response.json();
}

function BanksList({ banks }: { banks: Bank[] }) {
  return (
    <div>
      <h1>Banks List</h1>
      <div>
        {banks.map(bank => (
          <div key={bank.id}>{bank.name}</div>
        ))}
      </div>
    </div>
  );
}

export default async function Page() {
  let banks: Bank[] = [];

  try {
    banks = await fetchBanks();
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <BanksList banks={banks} />
      </Suspense>
    </>
  );
}
