'use client';

import { signOut } from 'next-auth/react';
import React, { useState } from 'react';
import { FaCog, FaSignOutAlt, FaUser, FaUserCircle } from 'react-icons/fa';

const UserMenu = () => {
  const [dropdownVisible, setDropdownVisible] = useState(false);

  const toggleDropdown = () => {
    setDropdownVisible(!dropdownVisible);
  };

  return (
    <div className='flex flex-col relative'>
      <FaUserCircle
        className='text-3xl cursor-pointer'
        onClick={toggleDropdown}
      />
      {dropdownVisible && (
        <div className='bg-gray-800 p-2 mt-2 absolute right-0 top-8 rounded shadow-lg flex flex-col gap-2'>
          <button className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'>
            <FaUser className='mr-2' /> Profile
          </button>
          <button className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'>
            <FaCog className='mr-2' /> Settings
          </button>
          <button
            className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'
            onClick={() => signOut()}
          >
            <FaSignOutAlt className='mr-2' /> Log out
          </button>
        </div>
      )}
    </div>
  );
};

export default UserMenu;
