import React from 'react';
import { MapPin, MessageCircle, UserPlus, Briefcase, GraduationCap, Globe, Star } from 'lucide-react';

export const MentorProfile = () => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
      <div className="flex flex-col md:flex-row items-start gap-6">
        <img
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=400&h=400&auto=format&fit=crop"
          alt="Andre Luis"
          className="w-32 h-32 rounded-full object-cover"
        />
        
        <div className="flex-1">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">Andre Luis Andrade Verde</h1>
              <p className="text-lg text-gray-600 dark:text-gray-300 mb-2">Senior Software Engineer at Google</p>
              <div className="flex items-center text-gray-500 dark:text-gray-400 mb-4">
                <MapPin size={16} className="mr-1" />
                <span>San Francisco, CA</span>
              </div>
            </div>
            
            <div className="flex gap-3">
              <button className="flex items-center gap-2 px-4 py-2 border border-blue-600 text-blue-600 rounded-lg hover:bg-blue-50 dark:border-blue-400 dark:text-blue-400">
                <UserPlus size={18} />
                <span>Follow</span>
              </button>
              <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
                <MessageCircle size={18} />
                <span>Message</span>
              </button>
            </div>
          </div>
          
          <div className="flex flex-wrap gap-2 mb-4">
            {['React', 'Node.js', 'System Design', 'Cloud Architecture'].map((skill) => (
              <span key={skill} className="px-3 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-full text-sm">
                {skill}
              </span>
            ))}
          </div>
          
          <div className="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
            <div className="flex items-center">
              <Star size={16} className="text-yellow-400 mr-1" />
              <span>4.9 (120 reviews)</span>
            </div>
            <div>•</div>
            <div>500+ sessions</div>
            <div>•</div>
            <div className="flex items-center">
              <Globe size={16} className="mr-1" />
              <span>English, Portuguese</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};