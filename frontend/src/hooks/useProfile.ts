import { useState } from 'react';

interface Profile {
  name: string;
  role: string;
  company: string;
  avatar: string;
  location: string;
  languages: string[];
  bio: string;
}

interface Stats {
  totalMentees: number;
  hoursSpent: number;
  avgRating: string;
  totalSessions: number;
}

interface Session {
  id: string;
  title: string;
  mentee: {
    name: string;
    avatar: string;
  };
  date: string;
  duration: string;
  platform: string;
}

const initialProfile: Profile = {
  name: 'Andre Luis',
  role: 'Senior Software Engineer',
  company: 'Google',
  avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=400&h=400&auto=format&fit=crop',
  location: 'San Francisco, CA',
  languages: ['English', 'Portuguese'],
  bio: 'Passionate about helping others grow in their tech careers. Specialized in system design and scalable architectures.'
};

const initialStats: Stats = {
  totalMentees: 45,
  hoursSpent: 120,
  avgRating: '4.9',
  totalSessions: 89
};

const initialSessions: Session[] = [
  {
    id: '1',
    title: 'System Design Discussion',
    mentee: {
      name: 'Sarah Chen',
      avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&h=150&auto=format&fit=crop'
    },
    date: 'Today, 3:00 PM',
    duration: '45 min',
    platform: 'Google Meet'
  },
  {
    id: '2',
    title: 'Career Growth Strategy',
    mentee: {
      name: 'Michael Park',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=150&h=150&auto=format&fit=crop'
    },
    date: 'Tomorrow, 2:00 PM',
    duration: '60 min',
    platform: 'Zoom'
  }
];

export const useProfile = () => {
  const [profile] = useState<Profile>(initialProfile);
  const [stats] = useState<Stats>(initialStats);
  const [upcomingSessions] = useState<Session[]>(initialSessions);

  return {
    profile,
    stats,
    upcomingSessions
  };
};