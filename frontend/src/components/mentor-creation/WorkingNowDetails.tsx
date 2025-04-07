import React, { useState } from 'react';
import { Linkedin, Instagram, Github, ChevronRight } from 'lucide-react';
import axios from 'axios';
import { API_URL } from '../../App';

interface WorkingNowDetailsProps {
  onNext: () => void;
}

export const WorkingNowDetails = ({ onNext }: WorkingNowDetailsProps) => {
  const [formData, setFormData] = useState({
    title: '',
    company: '',
    totalyear: '',
    month: '',
    linkedin: '',
    instagram : '',
    github: '',
  });

  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    const token = localStorage.getItem('token');
    if (!token) {
      setError('Please log in to continue');
      setIsLoading(false);
      return;
    }

    // Prepare working at payload
    const workingAtPayload = {
      title: formData.title,
      company: formData.company,
      totalyear: parseInt(formData.totalyear, 10),
      month: parseInt(formData.month, 10),
      linkedin: formData.linkedin,
      instagram: formData.instagram,
      github: formData.github,
    };

    try {
      // Post working at details
      const workingAtResponse = await axios.post(
        `${API_URL}/v1/workingat/create`,
        workingAtPayload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      if (workingAtResponse.status !== 201) {
        throw new Error('Failed to create working at details');
      }

      onNext();
    } catch (error: any) {
      setError(error.message || 'Error saving profile');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Current Title
        </label>
        <input
          type="text"
          required
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          placeholder="e.g., Senior Software Engineer"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Current Company
        </label>
        <input
          type="text"
          required
          value={formData.company}
          onChange={(e) => setFormData({ ...formData, company: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          placeholder="e.g., Google"
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Years of Experience
          </label>
          <input
            type="number"
            required
            min="0"
            max="50"
            value={formData.totalyear}
            onChange={(e) =>
              setFormData({ ...formData, totalyear: e.target.value })
            }
            className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Additional Months
          </label>
          <input
            type="number"
            required
            min="0"
            max="11"
            value={formData.month}
            onChange={(e) =>
              setFormData({ ...formData, month: e.target.value })
            }
            className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="space-y-4">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white">
          Social Media Links
        </h3>

        <div className="relative">
          <Linkedin
            size={20}
            className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
          />
          <input
            type="url"
            value={formData.linkedin}
            onChange={(e) =>
              setFormData({ ...formData, linkedin: e.target.value })
            }
            className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
            placeholder="LinkedIn Profile URL"
          />
        
        </div>

        <div className="relative">
          <Instagram
            size={20}
            className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
          />
          <input
            type="url"
            value={formData.instagram}
            onChange={(e) =>
              setFormData({ ...formData, instagram: e.target.value })
            }
            className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
            placeholder="Twitter Profile URL"
          />
        </div>

        <div className="relative">
          <Github
            size={20}
            className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
          />
          <input
            type="url"
            value={formData.github}
            onChange={(e) => setFormData({ ...formData, github: e.target.value })}
            className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
            placeholder="GitHub Profile URL"
          />
        </div>
      </div>

      {error && <div className="text-red-500 text-sm">{error}</div>}

      <div className="flex justify-end">
        <button
          type="submit"
          disabled={isLoading}
          className={`bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2 ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
        >
          {isLoading ? 'Submitting...' : 'Next'}
          {!isLoading && <ChevronRight size={18} />}
        </button>
      </div>
    </form>
  );
};