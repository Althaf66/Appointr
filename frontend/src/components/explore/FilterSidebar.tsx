import React from 'react';
import { Sliders, Clock, Briefcase, Code2 } from 'lucide-react';

const categories = [
  'Frontend Development',
  'Backend Development',
  'Mobile Development',
  'DevOps',
  'System Design',
  'Cloud Architecture',
];

const companies = [
  'Google',
  'Microsoft',
  'Amazon',
  'Meta',
  'Netflix',
  'Airbnb',
];

export const FilterSidebar = () => {
  return (
    <div className="w-64 bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
      <div className="flex items-center gap-2 mb-6">
        <Sliders size={20} className="text-blue-600 dark:text-blue-400" />
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Filters</h2>
      </div>

      <div className="space-y-6">
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Code2 size={18} className="text-gray-600 dark:text-gray-400" />
            <h3 className="font-medium text-gray-900 dark:text-white">Categories</h3>
          </div>
          {categories.map((category) => (
            <label key={category} className="flex items-center gap-2 mb-2">
              <input type="checkbox" className="rounded text-blue-600 dark:bg-gray-700 dark:border-gray-600" />
              <span className="text-sm text-gray-700 dark:text-gray-300">{category}</span>
            </label>
          ))}
        </div>

        <div>
          <div className="flex items-center gap-2 mb-3">
            <Briefcase size={18} className="text-gray-600 dark:text-gray-400" />
            <h3 className="font-medium text-gray-900 dark:text-white">Companies</h3>
          </div>
          {companies.map((company) => (
            <label key={company} className="flex items-center gap-2 mb-2">
              <input type="checkbox" className="rounded text-blue-600 dark:bg-gray-700 dark:border-gray-600" />
              <span className="text-sm text-gray-700 dark:text-gray-300">{company}</span>
            </label>
          ))}
        </div>

        <div>
          <div className="flex items-center gap-2 mb-3">
            <Clock size={18} className="text-gray-600 dark:text-gray-400" />
            <h3 className="font-medium text-gray-900 dark:text-white">Availability</h3>
          </div>
          {['Today', 'This Week', 'Next Week'].map((time) => (
            <label key={time} className="flex items-center gap-2 mb-2">
              <input type="checkbox" className="rounded text-blue-600 dark:bg-gray-700 dark:border-gray-600" />
              <span className="text-sm text-gray-700 dark:text-gray-300">{time}</span>
            </label>
          ))}
        </div>
      </div>
    </div>
  );
};