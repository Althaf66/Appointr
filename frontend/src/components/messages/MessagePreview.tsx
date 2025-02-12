import React from 'react';
import { Check } from 'lucide-react';
import type { Conversation } from '../../types/messages';

interface MessagePreviewProps {
  conversation: Conversation;
}

export const MessagePreview = ({ conversation }: MessagePreviewProps) => {
  const { mentor, lastMessage, unread } = conversation;

  return (
    <div className={`p-4 hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer ${
      unread ? 'bg-blue-50 dark:bg-blue-900/20' : ''
    }`}>
      <div className="flex items-center gap-3">
        <div className="relative">
          <img
            src={mentor.avatar}
            alt={mentor.name}
            className="w-12 h-12 rounded-full object-cover"
          />
          {mentor.online && (
            <div className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-white dark:border-gray-800" />
          )}
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex justify-between items-baseline">
            <h3 className="font-semibold text-gray-900 dark:text-white truncate">
              {mentor.name}
            </h3>
            <span className="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">
              {lastMessage.time}
            </span>
          </div>
          
          <div className="flex items-center gap-1">
            {lastMessage.sent && (
              <Check size={16} className="text-blue-600 dark:text-blue-400" />
            )}
            <p className={`text-sm truncate ${
              unread 
                ? 'text-gray-900 dark:text-white font-medium'
                : 'text-gray-500 dark:text-gray-400'
            }`}>
              {lastMessage.content}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};