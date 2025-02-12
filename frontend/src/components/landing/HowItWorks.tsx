import React from 'react';
import { Search, Calendar, Video, Star } from 'lucide-react';

const steps = [
  {
    icon: Search,
    title: 'Find Your Professionals',
    description: 'Browse through our curated list of expert professionals from worldwide.'
  },
  {
    icon: Calendar,
    title: 'Book a Session',
    description: 'Schedule a 1:1 session that fits your calendar and learning goals.'
  },
  {
    icon: Video,
    title: 'Meet Online',
    description: 'Connect with professionals virtually for personalized doubts and guidance.'
  },
  {
    icon: Star,
    title: 'Grow Together',
    description: 'Learn, solve doubts, and accelerate your professional growth.'
  }
];

export const HowItWorks = () => {
  return (
    <div className="py-20 bg-white dark:bg-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">How It Works</h2>
          <p className="text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
            Get started with mentorship in four simple steps. Our platform makes it easy to connect with the right mentor for your goals.
          </p>
        </div>
        
        <div className="grid md:grid-cols-4 gap-8">
          {steps.map((step, index) => (
            <div key={index} className="text-center">
              <div className="w-16 h-16 mx-auto mb-4 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center">
                <step.icon size={24} className="text-blue-600 dark:text-blue-400" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-gray-900 dark:text-white">{step.title}</h3>
              <p className="text-gray-600 dark:text-gray-400">{step.description}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};