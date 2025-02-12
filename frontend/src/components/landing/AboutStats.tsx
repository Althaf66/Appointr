import React from 'react';
import { Search, ArrowRight } from 'lucide-react';

export const AboutStats = () => {
  return (
    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-gray-900 dark:to-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center max-w-4xl mx-auto">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
            Finding an expert made simple
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
            Get personalized advice from people who have been there and done that.
          </p>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-8">
            Earn money for offering expert advice to those who need it.
          </p>
          
        </div>
      </div>
    </div>
  );
};