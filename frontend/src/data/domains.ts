export interface Domain {
  id: string;
  name: string;
  icon: string;
  description: string;
}

export const domains: Domain[] = [
  {
    id: 'engineering',
    name: 'Engineering',
    icon: 'Code2',
    description: 'Software Engineering, DevOps, and System Architecture'
  },
  {
    id: 'design',
    name: 'Design',
    icon: 'Palette',
    description: 'Product Design, UX/UI, and Visual Design'
  },
  {
    id: 'product',
    name: 'Product',
    icon: 'Layout',
    description: 'Product Management and Strategy'
  },
  {
    id: 'marketing',
    name: 'Marketing',
    icon: 'Target',
    description: 'Digital Marketing, Growth, and Analytics'
  },
  {
    id: 'data',
    name: 'Data',
    icon: 'Database',
    description: 'Data Science, Analytics, and Machine Learning'
  },
  {
    id: 'leadership',
    name: 'Leadership',
    icon: 'Users',
    description: 'Management, Leadership, and Career Development'
  }
];