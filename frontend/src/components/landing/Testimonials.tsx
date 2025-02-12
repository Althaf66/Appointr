import React from 'react';
import { Star } from 'lucide-react';

const testimonials = [
  {
    id: 1,
    content: "The mentorship I received was transformative. My mentor helped me navigate complex career decisions and provided invaluable industry insights.",
    author: "Sarah Chen",
    role: "Product Manager at Google",
    avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=150&h=150&auto=format&fit=crop"
  },
  {
    id: 2,
    content: "Appointr connected me with an amazing mentor who helped me transition into tech. The guidance was practical and actionable.",
    author: "Michael Rodriguez",
    role: "Software Engineer at Microsoft",
    avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?q=80&w=150&h=150&auto=format&fit=crop"
  },
  {
    id: 3,
    content: "The platform made it easy to find mentors who aligned with my career goals. The sessions were well-structured and incredibly helpful.",
    author: "Emily Wong",
    role: "UX Designer at Apple",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=150&h=150&auto=format&fit=crop"
  }
];

export const Testimonials = () => {
  return (
    <div className="py-20 bg-gray-50 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Success Stories</h2>
          <p className="text-gray-600 dark:text-gray-300">Hear from professionals who transformed their careers through mentorship</p>
        </div>
        
        <div className="grid md:grid-cols-3 gap-8">
          {testimonials.map((testimonial) => (
            <div key={testimonial.id} className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
              <div className="flex text-yellow-400 mb-4">
                {[...Array(5)].map((_, i) => (
                  <Star key={i} size={20} fill="currentColor" />
                ))}
              </div>
              <p className="text-gray-600 dark:text-gray-300 mb-6">{testimonial.content}</p>
              <div className="flex items-center">
                <img
                  src={testimonial.avatar}
                  alt={testimonial.author}
                  className="w-12 h-12 rounded-full mr-4"
                />
                <div>
                  <div className="font-semibold text-gray-900 dark:text-white">{testimonial.author}</div>
                  <div className="text-sm text-gray-500 dark:text-gray-400">{testimonial.role}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};