import React from 'react';
import { Users, Clock, Star, Calendar } from 'lucide-react';
import { useProfile } from '../../hooks/useProfile';

export const ProfileStats = () => {
  const { stats } = useProfile();

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {[
        { icon: Users, label: 'Total Mentees', value: stats.totalMentees },
        { icon: Clock, label: 'Hours Mentored', value: stats.hoursSpent },
        { icon: Star, label: 'Avg. Rating', value: stats.avgRating },
        { icon: Calendar, label: 'Sessions', value: stats.totalSessions }
      ].map((stat, index) => (
        <div key={index} className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-sm">
          <div className="flex items-center gap-2 mb-2">
            <stat.icon size={20} className="text-blue-600 dark:text-blue-400" />
            <span className="text-sm text-gray-600 dark:text-gray-400">{stat.label}</span>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">{stat.value}</p>
        </div>
      ))}
    </div>
  );
};