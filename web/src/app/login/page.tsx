// app/auth/signin/page.tsx
'use client';

import { signIn, useSession } from 'next-auth/react';

const SignIn = () => {
  // const { data: session } = useSession();
  // console.log(session);
  return (
    <div className='flex flex-col items-center justify-center min-h-screen py-2'>
      <h1 className='text-3xl font-bold mb-8'>Sign in</h1>
      <div className='flex flex-col space-y-4'>
        <div>
          <button
            onClick={() => signIn('google')}
            className='px-4 py-2 text-white bg-blue-500 rounded-md hover:bg-blue-600'
          >
            Sign in with Google
          </button>
        </div>
      </div>
    </div>
  );
};

export default SignIn;
