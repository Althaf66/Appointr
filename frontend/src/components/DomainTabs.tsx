import React from 'react';
import { domains } from '../data/domains';
import { Code2, Palette, Layout, Target, Database, Users } from 'lucide-react';

const iconMap = {
  Code2,
  Palette,
  Layout,
  Target,
  Database,
  Users,
};

export const DomainTabs = () => {
  return (
    <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex space-x-8 overflow-x-auto">
          {domains.map((domain) => {
            const Icon = iconMap[domain.icon as keyof typeof iconMap];
            return (
              <button
                key={domain.id}
                className="flex items-center space-x-2 py-4 px-1 border-b-2 border-transparent hover:border-blue-600 dark:hover:border-blue-400 text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white whitespace-nowrap"
              >
                <Icon size={20} />
                <span>{domain.name}</span>
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
};