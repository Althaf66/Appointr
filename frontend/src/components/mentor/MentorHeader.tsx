import React from 'react';
import { Search, Bell, User, MessageCircle, UserPlus } from 'lucide-react';

export const MentorHeader = () => {
  return (
    <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex-1 flex items-center">
            <span className="text-2xl font-bold text-blue-600 dark:text-blue-400 mr-8">Appointr</span>
            <div className="hidden md:block flex-1 max-w-2xl">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  type="text"
                  placeholder="Search mentors..."
                  className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <button className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full">
              <Bell size={20} />
            </button>
            <button className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full">
              <User size={20} />
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};