import React, { useEffect, useState, useRef } from "react";
import axios from "axios";
import { API_URL } from '../App';
import { MentorHeader } from "../components/mentor/MentorHeader";

const MessagesPage: React.FC = () => {
  const token = localStorage.getItem('token');
  const [conversations, setConversations] = useState<any[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<number | null>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const userId = token ? getUseridFromJWT(token) : 0;
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (token) {
      fetchConversations();
    }
  }, [token]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const fetchConversations = async () => {
    try {
      const res = await axios.get(`${API_URL}/messages/conversations`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setConversations(Array.isArray(res.data.data) ? res.data.data : []);
      console.log(res.data.data)
    } catch (error) {
      console.error("Error fetching conversations", error);
      setConversations([]);
    }
  };

  const fetchMessages = async (conversationID: number) => {
    try {
      const res = await axios.get(`${API_URL}/messages/${conversationID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      // Sort messages chronologically (oldest to newest)
      const sortedMessages = Array.isArray(res.data.data) 
        ? [...res.data.data].sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
        : [];
      setMessages(sortedMessages);
      setSelectedConversation(conversationID);
      console.log(res.data.data)
    } catch (error) {
      console.error("Error fetching messages", error);
      setMessages([]);
    }
  };

  function getUseridFromJWT(token: string): number | null {
    try {
      const payload = token.split('.')[1]; 
      const decodedPayload = atob(payload); 
      const jsonPayload = JSON.parse(decodedPayload); 
      return jsonPayload.sub || null;
    } catch (error) {
      console.error('Error decoding JWT token:', error);
      return null;
    }
  }

  const sendMessage = async () => {
    if (!selectedConversation || !newMessage.trim()) return;
    try {
      await axios.post(`${API_URL}/messages/conversations/${selectedConversation}/messages`,
        { content: newMessage },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setNewMessage("");
      fetchMessages(selectedConversation);
    } catch (error) {
      console.error("Error sending message", error);
    }
  };

  const getLastMessage = (conversation: any) => {
    return conversation.last_message.content || "No messages yet";
  };

  const getConversationPartner = (conversation: any) => {
    return conversation.other_user.username || "Unknown";
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const getInitials = (name: string) => {
    return name.split(' ').map(n => n[0]).join('').toUpperCase();
  };

  const formatMessageTime = (timestamp: string) => {
    if (!timestamp) return "00:00";
    try {
      const date = new Date(timestamp);
      return date.toTimeString().slice(0, 5);
    } catch (e) {
      return "00:00";
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <MentorHeader />
      <div className="flex h-screen">
        {/* Left sidebar - Conversations list */}
        <div className="w-1/3 border-r border-gray-800 flex flex-col">
          <div className="p-4 border-b border-gray-800">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Messages</h1>
            <div className="mt-4 relative">
              <input
                type="text"
                placeholder="Search messages"
                className="w-full bg-gray-800 rounded-md py-2 pl-10 pr-4 text-gray-300 focus:outline-none"
              />
              <svg className="absolute left-3 top-2.5 h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto">
            {conversations.length > 0 ? conversations.map((conv) => (
              <div 
                key={conv.id} 
                className={`flex items-center p-4 cursor-pointer hover:bg-gray-800 ${selectedConversation === conv.id ? 'bg-gray-800' : ''}`}
                onClick={() => fetchMessages(conv.id)}
              >
                <div className="flex-shrink-0 w-12 h-12 flex items-center justify-center bg-blue-600 rounded-full text-white font-bold">
                  {getInitials(getConversationPartner(conv))}
                </div>
                <div className="ml-3 flex-1 overflow-hidden">
                  <div className="flex justify-between items-center">
                    <p className="text-1xl font-bold text-gray-900 dark:text-white">{getConversationPartner(conv)}</p>
                    <p className="text-xs text-gray-500">{formatMessageTime(conv.last_message.created_at) || "00:00"}</p>
                  </div>
                  <p className="text-sm text-gray-400 truncate">{getLastMessage(conv)}</p>
                </div>
              </div>
            )) : (
              <p className="p-4 text-gray-500">No conversations found</p>
            )}
          </div>
        </div>

        {/* Right side - Messages */}
        <div className="w-2/3 flex flex-col">
          {selectedConversation ? (
            <>
              {/* Chat header */}
              <div className="flex items-center p-4 border-b border-gray-800">
                <div className="flex-shrink-0 w-10 h-10 flex items-center justify-center bg-blue-600 rounded-full text-white font-bold">
                  {getInitials(getConversationPartner(conversations.find(conv => conv.id === selectedConversation)))}
                </div>
                <div className="ml-3">
                  <p className="text-2xl font-bold text-gray-900 dark:text-white">{getConversationPartner(conversations.find(conv => conv.id === selectedConversation))}</p>
                  
                </div>
                <div className="ml-auto flex space-x-3">
                  <button className="text-gray-400 hover:text-white">
                    <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                    </svg>
                  </button>
                  <button className="text-gray-400 hover:text-white">
                    <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  </button>
                  <button className="text-gray-400 hover:text-white">
                    <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                </div>
              </div>

              {/* Messages - WhatsApp style (oldest at top, newest at bottom) */}
              <div className="flex-1 p-4 overflow-y-auto bg-gray-900 flex flex-col">
                {messages.length > 0 ? (
                  <div className="flex-1 flex flex-col justify-end">
                    <div className="space-y-4">
                      {messages.map((msg) => (
                        <div key={msg.id} className={`flex ${msg.sender.id === userId ? 'justify-end' : 'justify-start'}`}>
                          <div 
                            className={`max-w-xs p-3 rounded-2xl ${
                              msg.sender.id === userId 
                                ? 'bg-blue-600 text-white rounded-br-none' 
                                : 'bg-gray-800 text-white rounded-bl-none'
                            }`}
                          >
                            <p>{msg.content}</p>
                            <p className={`text-xs mt-1 ${msg.sender.id === userId ? 'text-blue-300' : 'text-gray-400'}`}>
                              {formatMessageTime(msg.created_at) || "00:00"}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                    <div ref={messagesEndRef}></div>
                  </div>
                ) : (
                  <div className="flex-1 flex items-center justify-center">
                    <p className="text-center text-gray-500">No messages yet</p>
                  </div>
                )}
              </div>

              {/* Message input */}
              <div className="p-4 border-t border-gray-800 flex items-center space-x-2">
                <button className="text-gray-400 hover:text-white">
                  <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </button>
                <button className="text-gray-400 hover:text-white">
                  <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                </button>
                <input
                  type="text"
                  placeholder="Type a message..."
                  className="flex-1 bg-gray-800 rounded-full px-4 py-2 text-white focus:outline-none"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  onKeyDown={handleKeyPress}
                />
                <button className="text-gray-400 hover:text-white">
                  <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
                <button 
                  onClick={sendMessage}
                  className="text-blue-400 hover:text-white"
                >
                  <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                </button>
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center">
              <p className="text-gray-500">Select a conversation to view messages</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MessagesPage;