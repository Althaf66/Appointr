import React from 'react';
import { Bell, Clock, Shield, Globe, Mail } from 'lucide-react';

export const ProfileSettings = () => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">Settings</h2>
      
      <div className="space-y-4">
        {[
          { icon: Bell, label: 'Notifications', description: 'Manage your notification preferences' },
          { icon: Clock, label: 'Availability', description: 'Set your mentoring schedule' },
          { icon: Shield, label: 'Privacy', description: 'Control your profile visibility' },
          { icon: Globe, label: 'Language', description: 'Change your preferred language' },
          { icon: Mail, label: 'Email Settings', description: 'Update email preferences' }
        ].map((setting, index) => (
          <button
            key={index}
            className="w-full flex items-center gap-4 p-4 text-left hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
              <setting.icon size={20} className="text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white">{setting.label}</h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">{setting.description}</p>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};