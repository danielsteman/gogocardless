import { FaUserCircle } from 'react-icons/fa';
import UserMenu from './components/UserMenu';

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
  let banks: Bank[] = [];

  try {
    banks = await fetchBanks();
  } catch (error) {
    console.error('Error fetching banks:', error);
  }

  return (
    <div className='flex'>
      <div className='ml-auto'>
        <UserMenu />
      </div>
    </div>
  );
}
