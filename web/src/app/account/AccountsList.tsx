'use client';

import { useRouter } from 'next/navigation';

type AccountsListProps = {
  accounts: string[];
};

const AccountsList: React.FC<AccountsListProps> = ({ accounts }) => {
  const router = useRouter();

  const handleAccountClick = (accountId: string) => {
    router.push(`/accounts/${accountId}`);
  };

  return (
    <div>
      <h1 className='text-2xl pb-8 font-bold'>Accounts</h1>
      {accounts.map((accountId, index) => (
        <div key={index} onClick={() => handleAccountClick(accountId)}>
          {accountId}
        </div>
      ))}
    </div>
  );
};

export default AccountsList;
