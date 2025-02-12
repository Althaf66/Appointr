import React from 'react';
import { MessageSidebar } from '../components/messages/MessageSidebar';
import { MessageThread } from '../components/messages/MessageThread';
import { MentorHeader } from '../components/mentor/MentorHeader';

export const MessagesPage = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <MentorHeader />
      <div className="max-w-7xl mx-auto">
        <div className="h-[calc(100vh-64px)] flex">
          <MessageSidebar />
          <MessageThread />
        </div>
      </div>
    </div>
  );
};