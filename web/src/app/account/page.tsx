'use client';

import { useEffect } from 'react';

const Callback: React.FC = () => {
  const agreementRef = localStorage.getItem('agreementRef');

  useEffect(() => {
    if (agreementRef) {
      console.log('agreementRef:', agreementRef);
    }
  }, [agreementRef]);

  return (
    <div>
      <h1 className='text-2xl pb-8 font-bold'>Callback Page</h1>
      <div>Agreement reference: {agreementRef}</div>
    </div>
  );
};

export default Callback;
