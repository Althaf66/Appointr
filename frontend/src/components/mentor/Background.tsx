import React from 'react';
import { Briefcase, GraduationCap } from 'lucide-react';

const experience = [
  {
    company: 'Google',
    role: 'Senior Software Engineer',
    duration: '2020 - Present',
    description: 'Leading development of cloud infrastructure projects and mentoring junior engineers.'
  },
  {
    company: 'Microsoft',
    role: 'Software Engineer',
    duration: '2017 - 2020',
    description: 'Worked on Azure cloud services and distributed systems.'
  }
];

const education = [
  {
    school: 'Stanford University',
    degree: 'Master of Science',
    field: 'Computer Science',
    year: '2017'
  },
  {
    school: 'University of California, Berkeley',
    degree: 'Bachelor of Science',
    field: 'Computer Science',
    year: '2015'
  }
];

export const Background = () => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-6">
      <div className="mb-8">
        <div className="flex items-center mb-4">
          <Briefcase className="text-blue-600 dark:text-blue-400 mr-2" size={24} />
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">Experience</h2>
        </div>
        
        <div className="space-y-6">
          {experience.map((exp, index) => (
            <div key={index} className="border-l-2 border-gray-200 dark:border-gray-700 pl-4">
              <h3 className="font-semibold text-gray-900 dark:text-white">{exp.role}</h3>
              <p className="text-gray-600 dark:text-gray-400">{exp.company}</p>
              <p className="text-sm text-gray-500 dark:text-gray-500">{exp.duration}</p>
              <p className="mt-2 text-gray-600 dark:text-gray-400">{exp.description}</p>
            </div>
          ))}
        </div>
      </div>
      
      <div>
        <div className="flex items-center mb-4">
          <GraduationCap className="text-blue-600 dark:text-blue-400 mr-2" size={24} />
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">Education</h2>
        </div>
        
        <div className="space-y-6">
          {education.map((edu, index) => (
            <div key={index} className="border-l-2 border-gray-200 dark:border-gray-700 pl-4">
              <h3 className="font-semibold text-gray-900 dark:text-white">{edu.school}</h3>
              <p className="text-gray-600 dark:text-gray-400">{edu.degree} in {edu.field}</p>
              <p className="text-sm text-gray-500 dark:text-gray-500">Class of {edu.year}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};