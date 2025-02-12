import React from 'react';
import { Phone, Video, Info } from 'lucide-react';
import type { Mentor } from '../../types/messages';

interface ThreadHeaderProps {
  mentor: Mentor;
}

export const ThreadHeader = ({ mentor }: ThreadHeaderProps) => {
  return (
    <div className="p-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <img
          src={mentor.avatar}
          alt={mentor.name}
          className="w-10 h-10 rounded-full object-cover"
        />
        <div>
          <h3 className="font-semibold text-gray-900 dark:text-white">{mentor.name}</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {mentor.online ? 'Active now' : 'Away'}
          </p>
        </div>
      </div>
      
      <div className="flex items-center gap-4">
        <button className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full">
          <Phone size={20} />
        </button>
        <button className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full">
          <Video size={20} />
        </button>
        <button className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full">
          <Info size={20} />
        </button>
      </div>
    </div>
  );
};