import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { ChevronRight, GalleryVerticalEnd, X } from 'lucide-react';
import { API_URL } from '../../App';

interface GigDetailsProps {
  onNext: () => void;
}

interface Expertise {
  name: string;
}

interface Discipline {
  id: number;
  field: string;
  subfield: string;
}

export const GigDetails = ({ onNext }: GigDetailsProps) => {
  const [formData, setFormData] = useState({
    title: '',
    amount: 0,
    description: '',
    expertise: '',
    discipline: [] as string[],
  });
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [expertiseOptions, setExpertiseOptions] = useState<string[]>([]);
  const [disciplineOptions, setDisciplineOptions] = useState<Discipline[]>([]);

  // Fetch expertise options
  useEffect(() => {
    const fetchExpertise = async () => {
      try {
        const response = await axios.get(`${API_URL}/expertise`);
        const expertiseData: Expertise[] = response.data.data;
        const expertiseNames = expertiseData.map(item => item.name);
        setExpertiseOptions(expertiseNames);
      } catch (err) {
        setError('Failed to fetch expertise options');
        console.error('Error fetching expertise:', err);
      }
    };

    fetchExpertise();
  }, []);

  // Fetch discipline options
  useEffect(() => {
    const fetchDisciplines = async () => {
      try {
        const response = await axios.get(`${API_URL}/discipline`);
        const disciplineData: { data: Discipline[] } = response.data;
        setDisciplineOptions(disciplineData.data);
      } catch (err) {
        setError('Failed to fetch discipline options');
        console.error('Error fetching disciplines:', err);
      }
    };

    fetchDisciplines();
  }, []);

  const handleDisciplineChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedOptions = Array.from(e.target.selectedOptions)
      .map(option => option.value)
      .filter(value => value !== '');
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

  // Filter discipline options based on selected expertise
  const filteredDisciplines = disciplineOptions.filter(
    discipline => discipline.field === formData.expertise
  );

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
          Amount ($)
        </label>
        <input
          type="number"
          required
          min="0"
          value={formData.amount}
          onChange={(e) => setFormData({ ...formData, amount: Number(e.target.value) })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          placeholder="e.g., 50"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Expertise Field
        </label>
        <select
          required
          value={formData.expertise}
          onChange={(e) => setFormData({ ...formData, expertise: e.target.value, discipline: [] })} // Reset discipline when expertise changes
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select your main expertise</option>
          {expertiseOptions.map((expertise) => (
            <option key={expertise} value={expertise}>
              {expertise}
            </option>
          ))}
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
          disabled={!formData.expertise} // Disable until expertise is selected
        >
          <option value="" disabled>
            {formData.expertise ? 'Select subcategories' : 'Select expertise first'}
          </option>
          {filteredDisciplines.map((discipline) => (
            <option key={discipline.id} value={discipline.subfield}>
              {discipline.subfield}
            </option>
          ))}
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