import React from 'react';
import { Search, ArrowRight } from 'lucide-react';

export const PlatformStats = () => {
  return (
    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-gray-900 dark:to-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center max-w-4xl mx-auto">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
            Unlock your earning potential
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
            Create your profile, set your price, and start sharing your expertise. Be in charge your time and earn money for your advice.
          </p>
          
          {/* <div className="max-w-2xl mx-auto mb-8">
            <div className="flex items-center bg-white dark:bg-gray-800 rounded-lg shadow-sm p-2">
              <Search className="text-gray-400 ml-2" size={20} />
              <input
                type="text"
                placeholder="Search by role, company, or expertise..."
                className="w-full px-4 py-2 bg-transparent focus:outline-none dark:text-white"
              />
              <button className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2">
                Search
                <ArrowRight size={18} />
              </button>
            </div>
          </div> */}
          
        </div>
      </div>
    </div>
  );
};