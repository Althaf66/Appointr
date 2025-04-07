import React from 'react';
import { MessageSquare, LogOut } from 'lucide-react';
import { ThemeToggle } from '../ThemeToggle';

function Logout() {
  localStorage.clear()
  return window.location.href = '/login'
}

export const MentorHeader = () => {
  return (
    <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <a href="/explore" className="text-2xl font-bold text-blue-600 dark:text-blue-400">Appointr</a>
          </div>
          <div className="flex items-center space-x-4">
            <ThemeToggle />
            <button
              onClick={() => window.location.href = '/messages'}
              className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full"
            >
              <MessageSquare size={20} />
            </button>
            <button
              onClick={Logout}
              className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full"
            >
              <LogOut size={20} />
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};