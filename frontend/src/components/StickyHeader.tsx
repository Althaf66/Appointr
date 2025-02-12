import React from 'react';
import { LogIn, UserPlus } from 'lucide-react';

export const StickyHeader = () => {
  return (
    <div className="sticky top-0 z-50 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-14">
          <div className="flex items-center space-x-4">
            {/* <span className="text-sm text-gray-500 dark:text-gray-400">Ready to grow?</span>   */}
            {/* <a href="#" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">Join 100k+ members</a> */}
          </div>
          <div className="flex items-center space-x-4">
            <button className="flex items-center space-x-1 text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">
              <LogIn size={18} />
              <span>Login</span>
            </button>
            <button className="flex items-center space-x-1 bg-blue-600 text-white px-4 py-1.5 rounded-lg hover:bg-blue-700">
              <UserPlus size={18} />
              <span>Sign Up</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};