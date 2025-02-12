export interface Mentor {
  id: string;
  name: string;
  avatar: string;
  online: boolean;
  role: string;
  company: string;
}

export interface Message {
  id: string;
  content: string;
  time: string;
  type: 'sent' | 'received';
  read: boolean;
}

export interface Conversation {
  id: string;
  mentor: Mentor;
  lastMessage: {
    content: string;
    time: string;
    sent: boolean;
  };
  messages: Message[];
  unread: boolean;
}

export interface MessagesState {
  conversations: Conversation[];
  activeThread: Conversation | null;
}