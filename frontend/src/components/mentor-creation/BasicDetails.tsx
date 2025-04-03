import React, { useState } from 'react';
import axios from 'axios';
import { ChevronRight } from 'lucide-react';

interface BasicDetailsProps {
  onNext: () => void;
}

const API_URL = 'http://localhost:8080/v1';

export const BasicDetails = ({ onNext }: BasicDetailsProps) => {
  const [formData, setFormData] = useState({
    Name: '',
    country: '',
    language: [] as string[], // Initialize as an array
  });
  const [languageInput, setLanguageInput] = useState(''); // Temporary input state

  const handleAddLanguage = () => {
    if (languageInput.trim()) {
      setFormData({
        ...formData,
        language: [...formData.language, languageInput.trim()],
      });
      setLanguageInput(''); // Clear input
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    try {
      const response = await axios.post(
        `${API_URL}/mentors/create`,
        formData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );
  
      if (response.status !== 200 && response.status !== 201) {
        throw new Error('Failed to create mentor');
      }
  
      onNext();
    } catch (error) {
      console.error('Error creating mentor:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Full Name
        </label>
        <input
          type="text"
          required
          value={formData.Name}
          onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          placeholder="Enter your full name"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Country
        </label>
        <select
          required
          value={formData.country}
          onChange={(e) => setFormData({ ...formData, country: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select your country</option>
          <option value="US">United States</option>
          <option value="UK">United Kingdom</option>
          <option value="CA">Canada</option>
          <option value="AU">Australia</option>
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Languages
        </label>
        <div className="flex gap-2">
          <input
            type="text"
            value={languageInput}
            onChange={(e) => setLanguageInput(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., English"
            onKeyPress={(e) => e.key === 'Enter' && handleAddLanguage()}
          />
          <button
            type="button"
            onClick={handleAddLanguage}
            className="bg-gray-200 px-4 py-2 rounded-lg hover:bg-gray-300"
          >
            Add
          </button>
        </div>
        <ul className="mt-2">
          {formData.language.map((lang, index) => (
            <li key={index} className="text-sm text-gray-700 dark:text-gray-300">
              {lang}
            </li>
          ))}
        </ul>
      </div>

      <div className="flex justify-end">
        <button
          type="submit"
          className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
        >
          Next
          <ChevronRight size={18} />
        </button>
      </div>
    </form>
  );
};