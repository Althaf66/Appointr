import { useState, useEffect } from 'react';
import type { MessagesState, Conversation, Message, Mentor } from '../types/messages';
import { API_URL } from '../App';

// Helper to get JWT token (adjust based on your auth setup)
const getToken = () => 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhcHBvaW50ciIsImV4cCI6MTc0MTA3MjAzMSwiaWF0IjoxNzQwODEyODMxLCJpc3MiOiJhcHBvaW50ciIsIm5iZiI6MTc0MDgxMjgzMSwic3ViIjoxfQ.sQQp7zP0BO7qV4kDWdhg1X1YFu2g579oakK05XLGFsI';

export const useMessages = () => {
  const [state, setState] = useState<MessagesState>({
    conversations: [],
    activeThread: null,
  });

  // Fetch conversations from backend
  const fetchConversations = async () => {
    try {
      const response = await fetch(`${API_URL}/messages/conversations`, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to fetch conversations: ${response.status} - ${errorText}`);
      }

      const responseData = await response.json();
      console.log('Conversations response:', responseData); // Debug log

      // Access the 'data' property, which contains the array
      const data = responseData.data;

      // Ensure data is an array
      if (!Array.isArray(data)) {
        throw new Error('Expected an array of conversations, but received: ' + JSON.stringify(data));
      }

      // Transform backend data to match frontend Conversation type
      const conversations = await Promise.all(
        data.map(async (conv: any) => {
          const messages = await fetchMessages(conv.id);
          const lastMessage = messages[messages.length - 1] || {
            id: '',
            content: '',
            time: '',
            type: 'received' as const,
            read: false,
          };
          // Use other_user.id instead of calculating from user1_id/user2_id
          const otherUserId = conv.other_user.id;
          const mentor = await fetchUser(otherUserId);

          return {
            id: conv.id.toString(),
            mentor,
            lastMessage: {
              content: conv.last_message.content,
              time: lastMessage.time || new Date(conv.last_message.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
              sent: conv.last_message.sender_id === getCurrentUserId(),
            },
            messages,
            unread: messages.some((msg: Message) => !msg.read && msg.type === 'received'),
          };
        })
      );

      setState((prev) => ({ ...prev, conversations }));
    } catch (error) {
      console.error('Error fetching conversations:', error);
    }
  };

  // Fetch messages for a conversation
  const fetchMessages = async (conversationId: string | number): Promise<Message[]> => {
    try {
      const response = await fetch(`${API_URL}/messages/${conversationId}?limit=50&offset=0`, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
      if (!response.ok) throw new Error('Failed to fetch messages');
      const responseData = await response.json();
      const data = responseData.data;
      console.log(`Messages for conversation ${conversationId}:`, data); // Debug log

      if (!Array.isArray(data)) {
        throw new Error('Expected an array of messages, but received: ' + JSON.stringify(data));
      }

      return data.map((msg: any) => ({
        id: msg.id.toString(),
        content: msg.content,
        time: new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        type: msg.sender_id === getCurrentUserId() ? 'sent' : 'received',
        read: msg.read,
      }));
    } catch (error) {
      console.error(`Error fetching messages for conversation ${conversationId}:`, error);
      return [];
    }
  };

  // Fetch user details (mentor)
  const fetchUser = async (userId: number): Promise<Mentor> => {
    try {
      const response = await fetch(`${API_URL}/users/${userId}`, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
      if (!response.ok) throw new Error('Failed to fetch user');
      const responseData = await response.json();
      const data = responseData.data;
      // console.log(`User ${userId} data:`, data); // Debug log

      return {
        id: data.id.toString(),
        name: data.username,
        avatar: `https://api.dicebear.com/7.x/initials/svg?seed=${data.username}`, // Replace with real avatar if available
        online: data.is_active, // Add logic if available
        role: data.role || 'Unknown', // Adjust if role exists in backend
        company: data.company || 'Unknown', // Adjust if company exists in backend
      };
    } catch (error) {
      console.error(`Error fetching user ${userId}:`, error);
      return {
        id: userId.toString(),
        name: `User ${userId}`,
        avatar: 'https://via.placeholder.com/150',
        online: false,
        role: 'Unknown',
        company: 'Unknown',
      };
    }
  };

  // Send a message
  const sendMessage = async (recipientId: string, content: string) => {
    try {
      const response = await fetch(`${API_URL}/messages/conversations/${recipientId}/messages`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${getToken()}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ recipient_id: parseInt(recipientId), content }),
      });
      if (!response.ok) throw new Error('Failed to send message');
      await fetchConversations(); // Refresh conversations
      if (state.activeThread) {
        setActiveThread(state.activeThread.id); // Refresh active thread
      }
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  // Mark conversation as read
  const markAsRead = async (conversationId: string) => {
    try {
      const response = await fetch(`${API_URL}/messages/${conversationId}/read`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
      if (!response.ok) throw new Error('Failed to mark as read');
      await fetchConversations(); // Refresh conversations
    } catch (error) {
      console.error('Error marking conversation as read:', error);
    }
  };

  // Set active thread and fetch its messages
  const setActiveThread = async (conversationId: string) => {
    const conversation = state.conversations.find((conv) => conv.id === conversationId);
    if (conversation) {
      const messages = await fetchMessages(conversationId);
      const updatedConversation = { ...conversation, messages };
      setState((prev) => ({ ...prev, activeThread: updatedConversation }));
      if (updatedConversation.unread) {
        await markAsRead(conversationId);
      }
    }
  };

  // Get unread count
  const getUnreadCount = async (): Promise<number> => {
    try {
      const response = await fetch(`${API_URL}/messages/unread`, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
      if (!response.ok) throw new Error('Failed to get unread count');
      const responseDatadata = await response.json();
      const data = responseDatadata.data;

      return data.count || 0;
    } catch (error) {
      console.error('Error fetching unread count:', error);
      return 0;
    }
  };

  // Placeholder for current user ID
  const getCurrentUserId = () => {
    return 1; // Replace with actual user ID from token/context
  };

  // Initial fetch
  useEffect(() => {
    fetchConversations();
  }, []);

  return {
    conversations: state.conversations,
    activeThread: state.activeThread,
    setActiveThread,
    sendMessage,
    getUnreadCount,
  };
};