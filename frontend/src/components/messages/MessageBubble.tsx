import React from 'react';
import { Check } from 'lucide-react';
import type { Message } from '../../types/messages';

interface MessageBubbleProps {
  message: Message;
}

export const MessageBubble = ({ message }: MessageBubbleProps) => {
  const isSent = message.type === 'sent';

  return (
    <div className={`flex ${isSent ? 'justify-end' : 'justify-start'}`}>
      <div className={`max-w-[70%] ${
        isSent 
          ? 'bg-blue-600 text-white rounded-t-2xl rounded-bl-2xl'
          : 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white rounded-t-2xl rounded-br-2xl'
      } px-4 py-2`}>
        <p>{message.content}</p>
        <div className={`flex items-center gap-1 text-xs mt-1 ${
          isSent ? 'text-blue-200' : 'text-gray-500 dark:text-gray-400'
        }`}>
          <span>{message.time}</span>
          {isSent && message.read && (
            <Check size={14} className="text-blue-200" />
          )}
        </div>
      </div>
    </div>
  );
};