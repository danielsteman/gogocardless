'use client';

import { useEffect, useState } from 'react';

const Account = () => {
  const [agreementRef, setAgreementRef] = useState<string | null>(null);
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    const ref = localStorage.getItem('agreementRef');
    setAgreementRef(ref);

    if (ref) {
      fetch(`http://localhost:3333/api/user/accounts?agreementRef=${ref}`)
        .then(response => response.json())
        .then(data => {
          setData(data);
        })
        .catch(error => {
          console.error('Error fetching data:', error);
        });
    }
  }, []);

  if (!agreementRef) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1 className='text-2xl pb-8 font-bold'>Callback Page</h1>
      <div>Agreement reference: {agreementRef}</div>
      {data && <pre>{JSON.stringify(data, null, 2)}</pre>}
    </div>
  );
};

export default Account;
