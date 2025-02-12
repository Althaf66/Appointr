import { Mentor } from '../types';

export const engineeringMentors: Mentor[] = [
  {
    id: '1',
    name: 'David Park',
    role: 'Staff Software Engineer',
    company: 'Netflix',
    avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['System Design', 'Backend Development', 'Microservices'],
    availability: 'This week',
  },
  {
    id: '2',
    name: 'Emily Chen',
    role: 'Engineering Manager',
    company: 'Meta',
    avatar: 'https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['Engineering Leadership', 'Career Development', 'Team Building'],
    availability: 'Next week',
  },
  {
    id: '3',
    name: 'Alex Kumar',
    role: 'Senior Frontend Engineer',
    company: 'Airbnb',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['React', 'Performance Optimization', 'Architecture'],
    availability: 'Tomorrow',
  },
];