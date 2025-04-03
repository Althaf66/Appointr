// import React, { useState } from 'react';
import { Search, Code2, ChevronDown } from 'lucide-react';
// import { domains } from '../../data/domains';

export const ExploreHeader = () => {
  // const [isDomainsOpen, setIsDomainsOpen] = useState(false);

  return (
    <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 transition-colors">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <div className="relative">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  type="text"
                  placeholder="Search by mentors name...."
                  className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              {/* <button
                onClick={() => setIsDomainsOpen(!isDomainsOpen)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              > */}
                {/* <ChevronDown size={20} /> */}
              {/* </button> */}

              {/* {isDomainsOpen && (
                <div className="absolute z-10 w-full mt-2 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700">
                  <div className="p-2">
                    {domains.map((domain) => (
                      <button
                        key={domain.id}
                        className="w-full text-left px-4 py-2 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-md flex items-center gap-3"
                        onClick={() => setIsDomainsOpen(false)}
                      >
                        <div className="p-1.5 bg-blue-50 dark:bg-blue-900/30 rounded-md">
                          <Code2 size={18} className="text-blue-600 dark:text-blue-400" />
                        </div>
                        <div>
                          <div className="font-medium text-gray-900 dark:text-white">{domain.name}</div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">{domain.description}</div>
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
              )} */}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};