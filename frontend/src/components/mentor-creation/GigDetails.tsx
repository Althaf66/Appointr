import React, { useState } from 'react';
import axios from 'axios';
import { ChevronRight, GalleryVerticalEnd, X } from 'lucide-react';
import { API_URL } from '../../App';

interface GigDetailsProps {
  onNext: () => void;
}

export const GigDetails = ({ onNext }: GigDetailsProps) => {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    expertise: '',
    discipline: [] as string[],
  });
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleDisciplineChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedOptions = Array.from(e.target.selectedOptions)
      .map(option => option.value)
      .filter(value => value !== ''); // Filter out the placeholder value
    const newDisciplines = [...new Set([...formData.discipline, ...selectedOptions])];
    setFormData({
      ...formData,
      discipline: newDisciplines,
    });
  };

  const handleRemoveDiscipline = (index: number) => {
    const updatedDiscipline = [...formData.discipline];
    updatedDiscipline.splice(index, 1);
    setFormData({
      ...formData,
      discipline: updatedDiscipline,
    });
  };

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

    try {
      const response = await axios.post(
        `${API_URL}/gigs/create`,
        formData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.status !== 201) {
        throw new Error('Failed to create gig');
      }

      onNext();
    } catch (error: any) {
      setError(error.message || 'Error creating gig');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Title
        </label>
        <input
          type="text"
          required
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          placeholder="e.g., Senior Software Engineer & Tech Mentor"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Bio
        </label>
        <textarea
          required
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500 h-32"
          placeholder="Tell us about yourself and your mentoring approach..."
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Expertise Field
        </label>
        <select
          required
          value={formData.expertise}
          onChange={(e) => setFormData({ ...formData, expertise: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select your main expertise</option>
          <option value="software">Software Development</option>
          <option value="design">Product Design</option>
          <option value="data">Data Science</option>
          <option value="product">Product Management</option>
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          <div className="flex items-center">
            <GalleryVerticalEnd size={16} className="mr-2" />
            Subcategory
          </div>
        </label>
        <select
          multiple
          value={formData.discipline}
          onChange={handleDisciplineChange}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500 h-32"
        >
          <option value="" disabled>Select subcategories</option>
          <option value="frontend">Frontend Development</option>
          <option value="backend">Backend Development</option>
          <option value="mobile">Mobile Development</option>
          <option value="devops">DevOps</option>
        </select>

        {formData.discipline.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-2">
            {formData.discipline.map((lang, index) => (
              <div 
                key={index} 
                className="inline-flex items-center bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-sm dark:bg-blue-900 dark:text-blue-200"
              >
                {lang}
                <button
                  type="button"
                  onClick={() => handleRemoveDiscipline(index)}
                  className="ml-2 text-blue-500 hover:text-blue-700 dark:text-blue-300 dark:hover:text-blue-100"
                >
                  <X size={14} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {error && <div className="text-red-500 text-sm">{error}</div>}

      <div className="flex justify-end">
        <button
          type="submit"
          disabled={isLoading}
          className={`bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2 ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
        >
          {isLoading ? 'Submitting...' : 'Submit'}
          {!isLoading && <ChevronRight size={18} />}
        </button>
      </div>
    </form>
  );
};