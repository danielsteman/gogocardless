'use client';

import React, { ReactNode } from 'react';
import UserMenu from '../components/UserMenu';

interface ProtectedLayoutProps {
  children: ReactNode;
}

const ProtectedLayout: React.FC<ProtectedLayoutProps> = ({ children }) => {
  return (
    <div className='min-h-screen flex flex-col'>
      <header className='bg-gray-900 text-white p-4 flex justify-between items-center'>
        <h1 className='text-lg ml-4'>PayPredict</h1>
        <UserMenu />
      </header>
      <main className='flex-grow p-8'>{children}</main>
      <footer className='bg-gray-900 text-white p-4 text-center'>
        &copy; {new Date().getFullYear()} PayPredict. All rights reserved.
      </footer>
    </div>
  );
};

export default ProtectedLayout;
