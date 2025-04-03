import React, { useState } from 'react';
import axios from 'axios';
import { Plus, Trash2, ChevronRight } from 'lucide-react';
import { API_URL } from '../../App';

interface EducationDetailsProps {
  onNext: () => void;
}

interface Education {
  id: string;
  year_from: string;
  year_to: string;
  degree: string;
  field: string;
  institute: string;
}

export const EducationDetails = ({ onNext }: EducationDetailsProps) => {
  const [educations, setEducations] = useState<Education[]>([
    { id: '1', year_from: '', year_to: '', degree: '', field: '', institute: '' },
  ]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const addEducation = () => {
    setEducations([
      ...educations,
      { id: Date.now().toString(), year_from: '', year_to: '', degree: '', field: '', institute: '' },
    ]);
  };

  const removeEducation = (id: string) => {
    if (educations.length > 1) {
      setEducations(educations.filter((edu) => edu.id !== id));
    }
  };

  const updateEducation = (id: string, field: keyof Education, value: string) => {
    setEducations(
      educations.map((edu) =>
        edu.id === id ? { ...edu, [field]: value } : edu
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
      // Send each education as a separate POST request
      for (const edu of educations) {
        const payload = {
          year_from: edu.year_from,
          year_to: edu.year_to,
          degree: edu.degree,
          field: edu.field,
          institute: edu.institute,
        };

        const response = await axios.post(
          `${API_URL}/education/create`,
          payload,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          }
        );

        if (response.status !== 201) {
          throw new Error(`Failed to create education: ${edu.degree} at ${edu.institute}`);
        }
      }

      onNext(); // Proceed to the next step on success
    } catch (error: any) {
      setError(error.message || 'Error creating education entries');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {educations.map((education, index) => (
        <div
          key={education.id}
          className="p-6 border border-gray-200 dark:border-gray-700 rounded-lg space-y-4"
        >
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Education {index + 1}
            </h3>
            {educations.length > 1 && (
              <button
                type="button"
                onClick={() => removeEducation(education.id)}
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
                min="1990"
                max={new Date().getFullYear()}
                value={education.year_from}
                onChange={(e) =>
                  updateEducation(education.id, 'year_from', e.target.value)
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
                min="1994"
                max={new Date().getFullYear()}
                value={education.year_to}
                onChange={(e) =>
                  updateEducation(education.id, 'year_to', e.target.value)
                }
                className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Degree
            </label>
            <input
              type="text"
              required
              value={education.degree}
              onChange={(e) =>
                updateEducation(education.id, 'degree', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Bachelor of Science"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Field of Study
            </label>
            <input
              type="text"
              required
              value={education.field}
              onChange={(e) =>
                updateEducation(education.id, 'field', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Computer Science"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Institute
            </label>
            <input
              type="text"
              required
              value={education.institute}
              onChange={(e) =>
                updateEducation(education.id, 'institute', e.target.value)
              }
              className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Stanford University"
            />
          </div>
        </div>
      ))}

      <button
        type="button"
        onClick={addEducation}
        className="w-full py-2 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-gray-600 dark:text-gray-400 hover:border-blue-500 hover:text-blue-500 flex items-center justify-center gap-2"
      >
        <Plus size={18} />
        Add Another Education
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