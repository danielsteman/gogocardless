'use client';

import { useEffect } from 'react';
import { useSearchParams } from 'next/navigation';

const Callback: React.FC = () => {
  const searchParams = useSearchParams();
  const agreementRef = searchParams.get('ref');

  useEffect(() => {
    if (agreementRef) {
      console.log('agreementRef:', agreementRef);
    }
  }, [agreementRef]);

  return (
    <div>
      <h1>Callback Page</h1>
    </div>
  );
};

export default Callback;
