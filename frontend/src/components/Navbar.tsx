import React, { useState, useEffect } from 'react';
import { Menu, MessageSquare, User } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';

export const Navbar = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Check for JWT token on component mount
  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
  }, []);

  return (
    <nav className="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 transition-colors">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <div className="flex items-center">
            <a href="/explore" className="text-2xl font-bold text-blue-600 dark:text-blue-400">Appointr</a>
          </div>
          
          <div className="hidden md:flex items-center space-x-8">
            <a href="/creatementor" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Become a Mentor</a>
            {/* <a href="#" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Community</a> */}
            {/* <a href="#" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Resources</a> */}
            <ThemeToggle />
            <button
              onClick={() => window.location.href = '/messages'}
              className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full"
            >
              <MessageSquare size={20} />
            </button>

            {isLoggedIn ? (
              // Show profile icon when logged in
              <a href="/profile" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">
                <User size={24} />
              </a>
            ) : (
              // Show login and signup buttons when not logged in
              <>
                <button className="bg-blue-600 dark:bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600">
                  <a href="/login">Login</a>
                </button>
                <button className="bg-grey-600 dark:bg-grey-500 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600">
                  <a href="/signup">Get Started</a>
                </button>
              </>
            )}
          </div>

          <div className="flex md:hidden">
            <ThemeToggle />
            {isLoggedIn ? (
              <a href="/profile" className="text-gray-700 dark:text-gray-300 ml-4">
                <User size={24} />
              </a>
            ) : (
              <button className="text-gray-700 dark:text-gray-300 ml-4">
                <Menu size={24} />
              </button>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};