export interface Domain {
  id: string;
  name: string;
  icon: string;
  description: string;
}

export const domains: Domain[] = [
  {
    id: 'engineering',
    name: 'IT',
    icon: 'Code2',
    description: 'Software Engineering, DevOps, and System Architecture'
  },
  {
    id: 'medical',
    name: 'Medical',
    icon: 'ShieldPlus',
    description: 'Product Design, UX/UI, and Visual Design'
  },
  {
    id: 'law',
    name: 'Law',
    icon: 'Scale',
    description: 'Product Management and Strategy'
  },
  {
    id: 'farming',
    name: 'Agriculture',
    icon: 'Tractor',
    description: 'Data Science, Analytics, and Machine Learning'
  },
  {
    id: 'marketing',
    name: 'Marketing',
    icon: 'Database',
    description: 'Digital Marketing, Growth, and Analytics'
  },
  {
    id: 'leadership',
    name: 'Leadership',
    icon: 'Users',
    description: 'Management, Leadership, and Career Development'
  }
];