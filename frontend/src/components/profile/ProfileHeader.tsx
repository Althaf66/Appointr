import React from 'react';
import { Camera, MapPin, Globe } from 'lucide-react';
import { useProfile } from '../../hooks/useProfile';

export const ProfileHeader = () => {
  const { profile } = useProfile();

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
      <div className="flex flex-col md:flex-row gap-6">
        <div className="relative">
          <img
            src={profile.avatar}
            alt={profile.name}
            className="w-32 h-32 rounded-full object-cover"
          />
          <button className="absolute bottom-0 right-0 p-2 bg-blue-600 text-white rounded-full hover:bg-blue-700">
            <Camera size={20} />
          </button>
        </div>
        
        <div className="flex-1">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{profile.name}</h1>
              <p className="text-gray-600 dark:text-gray-300">{profile.role} at {profile.company}</p>
            </div>
            <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              Edit Profile
            </button>
          </div>
          
          <div className="mt-4 flex flex-wrap gap-4 text-gray-600 dark:text-gray-300">
            <div className="flex items-center gap-1">
              <MapPin size={16} />
              <span>{profile.location}</span>
            </div>
            <div className="flex items-center gap-1">
              <Globe size={16} />
              <span>{profile.languages.join(', ')}</span>
            </div>
          </div>
          
          <p className="mt-4 text-gray-600 dark:text-gray-300">{profile.bio}</p>
        </div>
      </div>
    </div>
  );
};