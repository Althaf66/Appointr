import React from 'react';
import { Mail } from 'lucide-react';

export const Footer = () => {
  return (
    <footer className="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div className="max-w-4xl mx-auto px-4 py-6"> {/* Changed max-w-7xl to max-w-4xl and reduced padding */}
        <div className="flex flex-col items-center space-y-4">
          <div className="text-center">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Appointr</h3>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Connecting professionals worldwide
            </p>
          </div>
          <p className="text-gray-500 dark:text-gray-400 text-sm">
            © {new Date().getFullYear()} Appointr. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
};