import React, { useState } from 'react';
import { Send, Paperclip, Image, Smile } from 'lucide-react';
import { MessageBubble } from './MessageBubble';
import { ThreadHeader } from './ThreadHeader';
import { useMessages } from '../../hooks/useMessages';

export const MessageThread = () => {
  const [message, setMessage] = useState('');
  const { activeThread } = useMessages();

  const handleSend = () => {
    if (message.trim()) {
      // Handle sending message
      setMessage('');
    }
  };

  if (!activeThread) {
    return (
      <div className="flex-1 flex items-center justify-center bg-white dark:bg-gray-800">
        <div className="text-center">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            Select a conversation
          </h3>
          <p className="text-gray-500 dark:text-gray-400">
            Choose a mentor to start messaging
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-white dark:bg-gray-800">
      <ThreadHeader mentor={activeThread.mentor} />
      
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {activeThread.messages.map((msg) => (
          <MessageBubble key={msg.id} message={msg} />
        ))}
      </div>
      
      <div className="p-4 border-t border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-2">
          <button className="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300">
            <Image size={20} />
          </button>
          <button className="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300">
            <Paperclip size={20} />
          </button>
          
          <input
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type a message..."
            className="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
            onKeyPress={(e) => e.key === 'Enter' && handleSend()}
          />
          
          <button className="p-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300">
            <Smile size={20} />
          </button>
          <button
            onClick={handleSend}
            disabled={!message.trim()}
            className="p-2 text-blue-600 dark:text-blue-400 hover:text-blue-700 disabled:opacity-50"
          >
            <Send size={20} />
          </button>
        </div>
      </div>
    </div>
  );
};