import React from 'react';
import { ArrowRight } from 'lucide-react';

const similarMentors = [
  {
    id: '1',
    name: 'Sarah Chen',
    role: 'Staff Engineer',
    company: 'Meta',
    avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&h=150&auto=format&fit=crop'
  },
  {
    id: '2',
    name: 'Michael Park',
    role: 'Engineering Manager',
    company: 'Amazon',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=150&h=150&auto=format&fit=crop'
  }
];

export const SimilarMentors = () => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white">Similar Mentors</h2>
        <button className="text-blue-600 dark:text-blue-400 flex items-center gap-1 hover:underline">
          View all
          <ArrowRight size={16} />
        </button>
      </div>
      
      <div className="grid md:grid-cols-2 gap-4">
        {similarMentors.map((mentor) => (
          <div key={mentor.id} className="flex items-center gap-4 p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
            <img
              src={mentor.avatar}
              alt={mentor.name}
              className="w-12 h-12 rounded-full object-cover"
            />
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white">{mentor.name}</h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">{mentor.role} at {mentor.company}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};