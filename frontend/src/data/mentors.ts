import { Mentor } from '../types';

export const mentors: Mentor[] = [
  {
    id: '1',
    name: 'Sarah Chen',
    role: 'Senior Product Designer',
    company: 'Google',
    avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['UX Design', 'Product Strategy', 'Design Systems'],
    availability: 'Next week',
  },
  {
    id: '2',
    name: 'Michael Rodriguez',
    role: 'Engineering Manager',
    company: 'Microsoft',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['Leadership', 'Engineering', 'Career Growth'],
    availability: 'This week',
  },
  {
    id: '3',
    name: 'Priya Sharma',
    role: 'Product Manager',
    company: 'Amazon',
    avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=150&h=150&auto=format&fit=crop',
    expertise: ['Product Management', 'Strategy', 'Growth'],
    availability: 'Tomorrow',
  },
];