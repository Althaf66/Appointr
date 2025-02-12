import { useState } from 'react';
import type { MessagesState, Conversation } from '../types/messages';

const initialState: MessagesState = {
  conversations: [
    {
      id: '1',
      mentor: {
        id: '1',
        name: 'Andre Luis',
        avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=150&h=150&auto=format&fit=crop',
        online: true,
        role: 'Senior Software Engineer',
        company: 'Google'
      },
      lastMessage: {
        content: 'Looking forward to our session!',
        time: '2:30 PM',
        sent: true
      },
      messages: [
        {
          id: '1',
          content: 'Hi Andre, I would love to learn more about system design.',
          time: '2:15 PM',
          type: 'sent',
          read: true
        },
        {
          id: '2',
          content: 'Looking forward to our session!',
          time: '2:30 PM',
          type: 'received',
          read: true
        }
      ],
      unread: false
    },
    {
      id: '2',
      mentor: {
        id: '2',
        name: 'Sarah Chen',
        avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&h=150&auto=format&fit=crop',
        online: false,
        role: 'Staff Engineer',
        company: 'Meta'
      },
      lastMessage: {
        content: 'Let me know if you have any questions!',
        time: '9:45 AM',
        sent: false
      },
      messages: [
        {
          id: '1',
          content: 'Let me know if you have any questions!',
          time: '9:45 AM',
          type: 'received',
          read: false
        }
      ],
      unread: true
    }
  ],
  activeThread: null
};

export const useMessages = () => {
  const [state, setState] = useState<MessagesState>(initialState);

  const setActiveThread = (conversation: Conversation) => {
    setState((prev) => ({
      ...prev,
      activeThread: conversation
    }));
  };

  return {
    conversations: state.conversations,
    activeThread: state.activeThread,
    setActiveThread
  };
};