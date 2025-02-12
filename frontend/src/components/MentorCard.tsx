import React from 'react';
import { Calendar, Clock } from 'lucide-react';
import type { Mentor } from '../types';

interface MentorCardProps {
  mentor: Mentor;
}

export const MentorCard = ({ mentor }: MentorCardProps) => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm hover:shadow-md transition-shadow p-6 border border-gray-200 dark:border-gray-700">
      <div className="flex items-start space-x-4">
        <img
          src={mentor.avatar}
          alt={mentor.name}
          className="w-16 h-16 rounded-full object-cover"
        />
        <div className="flex-1">
          <h3 className="font-semibold text-lg text-gray-900 dark:text-white">{mentor.name}</h3>
          <p className="text-gray-600 dark:text-gray-400">{mentor.role} at {mentor.company}</p>
          
          <div className="mt-4 flex flex-wrap gap-2">
            {mentor.expertise.map((skill, index) => (
              <span
                key={index}
                className="bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 text-sm px-3 py-1 rounded-full"
              >
                {skill}
              </span>
            ))}
          </div>
          
          <div className="mt-4 flex items-center text-sm text-gray-500 dark:text-gray-400">
            <Calendar size={16} className="mr-1" />
            <span>Available {mentor.availability}</span>
          </div>
        </div>
      </div>
      
      <div className="mt-4 pt-4 border-t dark:border-gray-700">
        <button className="w-full bg-white dark:bg-gray-800 border border-blue-600 dark:border-blue-400 text-blue-600 dark:text-blue-400 px-4 py-2 rounded-lg hover:bg-blue-50 dark:hover:bg-blue-900/30 transition-colors">
          Book Session
        </button>
      </div>
    </div>
  );
};