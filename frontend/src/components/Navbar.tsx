import React from 'react';
import { Search, Menu } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';

export const Navbar = () => {
  return (
    <nav className="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 transition-colors">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <div className="flex items-center">
            <a href="/" className="text-2xl font-bold text-blue-600 dark:text-blue-400">Appointr</a>
          </div>
          
          <div className="hidden md:flex items-center space-x-8">
            <a href="#" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Find a Mentor</a>
            {/* <a href="#" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Community</a> */}
            {/* <a href="#" className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400">Resources</a> */}
            <ThemeToggle />
            <button className="bg-blue-600 dark:bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600">
              <a href="/login">Login</a>
            </button>
            <button className="bg-grey-600 dark:bg-grey-500 text-white px-4 py-2 rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600">
              <a href="/signup">Get Started</a>
            </button>
          </div>

          <div className="flex md:hidden">
            <ThemeToggle />
            <button className="text-gray-700 dark:text-gray-300 ml-4">
              <Menu size={24} />
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};