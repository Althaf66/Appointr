import React from 'react';
import { Calendar, Clock, Video } from 'lucide-react';
import { useProfile } from '../../hooks/useProfile';

export const UpcomingSessions = () => {
  const { upcomingSessions } = useProfile();

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">Upcoming Sessions</h2>
      
      <div className="space-y-4">
        {upcomingSessions.map((session) => (
          <div key={session.id} className="flex items-center gap-4 p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div className="flex-shrink-0 w-16 h-16">
              <img
                src={session.mentee.avatar}
                alt={session.mentee.name}
                className="w-full h-full rounded-full object-cover"
              />
            </div>
            
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-gray-900 dark:text-white">{session.title}</h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">{session.mentee.name}</p>
              
              <div className="mt-2 flex flex-wrap gap-4 text-sm text-gray-500 dark:text-gray-400">
                <div className="flex items-center gap-1">
                  <Calendar size={16} />
                  <span>{session.date}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Clock size={16} />
                  <span>{session.duration}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Video size={16} />
                  <span>{session.platform}</span>
                </div>
              </div>
            </div>
            
            <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              Join
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};