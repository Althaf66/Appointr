import React, { useState } from 'react';
import { Calendar, Clock, Users } from 'lucide-react';

const sessions = [
  {
    id: '1',
    title: '1:1 Career Guidance Session',
    duration: '45 mins',
    spots: 3,
    price: 'Free',
    date: '2024-03-20'
  },
  // {
  //   id: '2',
  //   title: 'Technical Interview Prep',
  //   duration: '60 mins',
  //   spots: 2,
  //   price: 'Free',
  //   date: '2024-03-21'
  // }
];

export const BookingSection = () => {
  const [selectedSession, setSelectedSession] = useState(sessions[0]);

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
      <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-6">Book a Session</h2>
      
      <div className="grid md:grid-cols-2 gap-6">
        {sessions.map((session) => (
          <div
            key={session.id}
            className={`p-4 border-2 rounded-lg cursor-pointer transition-colors ${
              selectedSession.id === session.id
                ? 'border-blue-600 dark:border-blue-400'
                : 'border-gray-200 dark:border-gray-700'
            }`}
            onClick={() => setSelectedSession(session)}
          >
            <h3 className="font-semibold text-gray-900 dark:text-white mb-3">{session.title}</h3>
            <div className="flex items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
              <div className="flex items-center">
                <Clock size={16} className="mr-1" />
                <span>{session.duration}</span>
              </div>
              <div className="flex items-center">
                <Users size={16} className="mr-1" />
                <span>{session.spots} spots left</span>
              </div>
              <div className="flex items-center">
                <Calendar size={16} className="mr-1" />
                <span>{new Date(session.date).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
      
      <button className="w-full mt-6 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700">
        Book Session
      </button>
    </div>
  );
};