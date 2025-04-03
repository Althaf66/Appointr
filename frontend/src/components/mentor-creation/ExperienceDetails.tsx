import React, { useState } from 'react';
import axios from 'axios';
import { Plus, Trash2, ChevronRight } from 'lucide-react';
import { API_URL } from '../../App';

interface ExperienceDetailsProps {
  onNext: () => void;
}

interface Experience {
  id: string;
  year_from: string;
  year_to: string;
  title: string;
  company: string;
  description: string;
}

export const ExperienceDetails = ({ onNext }: ExperienceDetailsProps) => {
  const [experiences, setExperiences] = useState<Experience[]>([
    { id: '1', year_from: '', year_to: '', title: '', company: '', description: '' },
  ]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const addExperience = () => {
    setExperiences([
      ...experiences,
      { id: Date.now().toString(), year_from: '', year_to: '', title: '', company: '', description: '' },
    ]);
  };

  const removeExperience = (id: string) => {
    if (experiences.length > 1) {
      setExperiences(experiences.filter((exp) => exp.id !== id));
    }
  };

  const updateExperience = (id: string, field: keyof Experience, value: string) => {
    setExperiences(
      experiences.map((exp) =>
        exp.id === id ? { ...exp, [field]: value } : exp
      )
    );
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
      // Send each experience as a separate POST request
      for (const exp of experiences) {
        const payload = {
          year_from: exp.year_from,
          year_to: exp.year_to,
          title: exp.title,
          company: exp.company,
          description: exp.description,
        };
        console.log(payload)
        const response = await axios.post(
          `${API_URL}/experience/create`,
          payload,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          }
        );

        if (response.status !== 201) {
          throw new Error(`Failed to create experience: ${exp.title}`);
        }
      }

      onNext(); // Proceed to the next step on success
    } catch (error: any) {
      setError(error.message || 'Error creating experiences');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {experiences.map((experience, index) => (
        <div
          key={experience.id}
          className="p-6 border border-gray-200 dark:border-gray-700 rounded-lg space-y-4"
        >
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Experience {index + 1}
            </h3>
            {experiences.length > 1 && (
              <button
                type="button"
                onClick={() => removeExperience(experience.id)}
                className="text-red-500 hover:text-red-600"
              >
                <Trash2 size={18} />
              </button>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Year From
              </label>
              <input
                type="number"
                required
                min="2000"
                max={new Date().getFullYear()}
                value={experience.year_from}
                onChange={(e) =>
                  updateExperience(experience.id, 'year_from', e.target.value)
                }
                className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Year To
              </label>
              <input
                type="number"
                required
                min="2001"
                max={new Date().getFullYear()}
                value={experience.year_to}
                onChange={(e) =>
                  updateExperience(experience.id, 'year_to', e.target.value)
                }
                className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Title
            </label>
            <input
              type="text"
              required
              value={experience.title}
              onChange={(e) =>
                updateExperience(experience.id, 'title', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Senior Software Engineer"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Company
            </label>
            <input
              type="text"
              required
              value={experience.company}
              onChange={(e) =>
                updateExperience(experience.id, 'company', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Google"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Description
            </label>
            <input
              type="text"
              required
              value={experience.description}
              onChange={(e) =>
                updateExperience(experience.id, 'description', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="Describe your role and responsibilities"
            />
          </div>
        </div>
      ))}

      <button
        type="button"
        onClick={addExperience}
        className="w-full py-2 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-gray-600 dark:text-gray-400 hover:border-blue-500 hover:text-blue-500 flex items-center justify-center gap-2"
      >
        <Plus size={18} />
        Add Another Experience
      </button>

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