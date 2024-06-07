'use client';

import React, { useState } from 'react';
import { FaCog, FaSignOutAlt, FaUser, FaUserCircle } from 'react-icons/fa';

const UserMenu = () => {
  const [dropdownVisible, setDropdownVisible] = useState(false);

  const toggleDropdown = () => {
    setDropdownVisible(!dropdownVisible);
  };

  return (
    <div className='flex flex-col w-full p-8 relative'>
      <FaUserCircle
        className='text-3xl cursor-pointer'
        onClick={toggleDropdown}
      />
      {dropdownVisible && (
        <div className='bg-gray-800 p-2 mt-2 absolute right-8 top-16 rounded shadow-lg flex flex-col gap-2'>
          <button className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'>
            <FaUser className='mr-2' /> Profile
          </button>
          <button className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'>
            <FaCog className='mr-2' /> Settings
          </button>
          <button className='flex items-center w-full text-left text-white py-2 px-4 rounded hover:bg-gray-700'>
            <FaSignOutAlt className='mr-2' /> Log out
          </button>
        </div>
      )}
    </div>
  );
};

export default UserMenu;
