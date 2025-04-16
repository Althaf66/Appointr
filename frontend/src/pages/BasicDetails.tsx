import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Loader2, ChevronRight, User, MapPin, Languages, X, Plus } from 'lucide-react';
import { API_URL } from '../App';

interface Country {
  countryname: string;
}

export const BasicDetailsPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    Name: '',
    country: '',
    language: [] as string[],
  });
  const [languageInput, setLanguageInput] = useState('');
  const [countries, setCountries] = useState<Country[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchCountries = async () => {
      setIsLoading(true);
      try {
        const response = await axios.get(`${API_URL}/v1/countries`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });
        setCountries(response.data.data);
        setIsLoading(false);
      } catch (err) {
        console.error('Error fetching countries:', err);
        setError('Failed to load countries');
        setIsLoading(false);
      }
    };
    fetchCountries();
  }, []);

  const handleAddLanguage = () => {
    if (languageInput.trim() && !formData.language.includes(languageInput.trim())) {
      setFormData({
        ...formData,
        language: [...formData.language, languageInput.trim()],
      });
      setLanguageInput('');
    }
  };

  const handleRemoveLanguage = (index: number) => {
    const updatedLanguages = [...formData.language];
    updatedLanguages.splice(index, 1);
    setFormData({
      ...formData,
      language: updatedLanguages,
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
        `${API_URL}/v1/mentors/create`,
        formData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.status !== 201) {
        throw new Error('Failed to create mentor profile');
      }

      navigate('/creategig');
    } catch (error: any) {
      setError(error.response?.data?.message || 'Error creating mentor profile');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
            Create Your Mentor Profile
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
            Let's start with some basic information about you
          </p>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8">
          <div className="flex items-center mb-6">
            <h3 className="text-xl font-medium text-gray-900 dark:text-white">Basic Details</h3>
          </div>
          
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  <div className="flex items-center">
                    <User size={16} className="mr-2" />
                    Full Name
                  </div>
                </label>
                <input
                  type="text"
                  required
                  value={formData.Name}
                  onChange={(e) => setFormData({ ...formData, Name: e.target.value })}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  placeholder="Enter your full name"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  <div className="flex items-center">
                    <MapPin size={16} className="mr-2" />
                    Country
                  </div>
                </label>
                <div className="relative">
                <select
          required
          value={formData.country}
          onChange={(e) => setFormData({ ...formData, country: e.target.value })}
          className="w-full px-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
          disabled={isLoading} // Disable while loading
        >
                    <option value="">Select your country</option>
          {countries.map((country) => (
            <option key={country.countryname} value={country.countryname}>
              {country.countryname}
            </option>
          ))}
        </select>
        {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
                  
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  <div className="flex items-center">
                    <Languages size={16} className="mr-2" />
                    Languages
                  </div>
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={languageInput}
                    onChange={(e) => setLanguageInput(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                    placeholder="e.g., English, Spanish, French"
                    onKeyPress={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleAddLanguage();
                      }
                    }}
                  />
                  <button
                    type="button"
                    onClick={handleAddLanguage}
                    className="flex items-center justify-center px-4 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors dark:bg-gray-600 dark:hover:bg-gray-500 dark:text-white"
                  >
                    <Plus size={20} />
                  </button>
                </div>
                
                {formData.language.length > 0 && (
                  <div className="mt-3 flex flex-wrap gap-2">
                    {formData.language.map((lang, index) => (
                      <div 
                        key={index} 
                        className="inline-flex items-center bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-sm dark:bg-blue-900 dark:text-blue-200"
                      >
                        {lang}
                        <button
                          type="button"
                          onClick={() => handleRemoveLanguage(index)}
                          className="ml-2 text-blue-500 hover:text-blue-700 dark:text-blue-300 dark:hover:text-blue-100"
                        >
                          <X size={14} />
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg dark:bg-red-900/30 dark:border-red-800 dark:text-red-400">
                {error}
              </div>
            )}

            <div className="pt-4">
              <button
                type="submit"
                disabled={isLoading}
                className="w-full flex justify-center items-center py-3 px-6 border border-transparent rounded-lg shadow-sm text-base font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed dark:bg-blue-700 dark:hover:bg-blue-600"
              >
                {isLoading ? (
                  <>
                    <Loader2 size={20} className="animate-spin mr-2" />
                    Creating Profile...
                  </>
                ) : (
                  <>
                    Create Gig
                    <ChevronRight size={20} className="ml-2" />
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
        
      </div>
    </div>
  );
};