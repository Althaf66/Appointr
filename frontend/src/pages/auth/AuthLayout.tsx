import React from 'react';
import { Code2 } from 'lucide-react';
import { ThemeToggle } from '../../components/ThemeToggle';

interface AuthLayoutProps {
  children: React.ReactNode;
}

export const AuthLayout = ({ children }: AuthLayoutProps) => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex">
      {/* Left side - Form */}
      <div className="flex-1 flex flex-col justify-center items-center px-4 sm:px-6 lg:px-8 py-12">
        <div className="flex items-center mb-8">
          <Code2 size={32} className="text-blue-600 dark:text-blue-400" />
          <span className="ml-2 text-2xl font-bold text-gray-900 dark:text-white">Appointr</span>
        </div>
        {children}
      </div>

      {/* Right side - Image */}
      <div className="hidden lg:block relative flex-1">
        <div className="absolute top-4 right-4">
          <ThemeToggle />
        </div>
        <img
          src="https://images.unsplash.com/photo-1522071820081-009f0129c71c?q=80&w=2070&auto=format&fit=crop"
          alt="People collaborating"
          className="w-full h-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-gray-900 via-gray-900/40 to-transparent">
          <div className="absolute bottom-0 left-0 right-0 p-8 text-white">
            <blockquote className="text-xl font-medium mb-2">
              "The mentorship I received through Appointr transformed my career path completely."
            </blockquote>
            <p className="text-gray-200">Sarah Chen, Product Manager at Google</p>
          </div>
        </div>
      </div>
    </div>
  );
};