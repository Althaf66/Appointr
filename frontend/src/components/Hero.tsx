import React from 'react';
import { Search } from 'lucide-react';

export const Hero = () => {
  return (
    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 py-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
            Learn from the World's<br />Top Professionals
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            Book 1:1 mentoring sessions with leaders from the most innovative companies
          </p>
          
          <div className="max-w-2xl mx-auto">
            <div className="flex items-center bg-white rounded-lg shadow-sm p-2">
              <Search className="text-gray-400 ml-2" size={20} />
              <input
                type="text"
                placeholder="Search by role, company, or expertise..."
                className="w-full px-4 py-2 focus:outline-none"
              />
              <button className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700">
                Search
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};