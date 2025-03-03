// components/messages/MessageSidebar.tsx
import React from 'react';
import { Search, Edit } from 'lucide-react';
import { MessagePreview } from './MessagePreview';
import { useMessages } from '../../hooks/useMessages';

export const MessageSidebar = () => {
  const { conversations, setActiveThread } = useMessages();

  return (
    <div className="w-96 border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Messages</h2>
          <button className="text-blue-600 dark:text-blue-400 hover:text-blue-700">
            <Edit size={20} />
          </button>
        </div>
        
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
          <input
            type="text"
            placeholder="Search messages"
            className="w-full pl-10 pr-4 py-2 bg-gray-100 dark:bg-gray-700 border-0 rounded-lg focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="overflow-y-auto h-[calc(100vh-180px)]">
        {conversations.map((conversation) => (
          <div
            key={conversation.id}
            onClick={() => {
              console.log('Selecting conversation:', conversation.id);
              setActiveThread(conversation.id);
            }}
            className="cursor-pointer"
          >
            <MessagePreview conversation={conversation} />
          </div>
        ))}
      </div>
    </div>
  );
};