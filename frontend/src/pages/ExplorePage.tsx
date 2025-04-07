import { useState, useEffect } from 'react';
import { Navbar } from '../components/Navbar';
import { DomainTabs } from '../components/DomainTabs';
import { ExploreHeader } from '../components/explore/ExploreHeader';
import { MentorCard } from '../components/MentorCard';
import { Footer } from '../components/Footer';
import axios from 'axios';
import { Mentor } from '../types';

export const ExplorePage = () => {
  const [mentors, setMentors] = useState<Mentor[]>([]);
  const [loading, setLoading] = useState(false); // Changed to false since no initial loading
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [selectedSubfield, setSelectedSubfield] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Get auth token from localStorage
  const getAuthToken = () => {
    return localStorage.getItem('token');
  };

  // Configure axios with authentication headers
  const getAxiosConfig = () => {
    const token = getAuthToken();
    return {
      headers: {
        'Authorization': token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json'
      }
    };
  };

  // Function to update mentors from DomainTabs component
  const updateMentors = (newMentors: Mentor[], domain: string | null, subfield: string | null) => {
    setMentors(newMentors);
    setSelectedDomain(domain);
    setSelectedSubfield(subfield);
    setError(null);
    setLoading(false);
  };

  // Check if token exists and redirect to login if not
  useEffect(() => {
    const token = getAuthToken();
    if (!token) {
      setError('Authentication required. Please log in to view mentors.');
    }
  }, []);

  return (
    <div className="flex flex-col h-screen">
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
      <Navbar />
      <ExploreHeader />
      <DomainTabs 
        onMentorsUpdate={updateMentors} 
      />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex gap-8">
          <div className="flex-1">
            <div className="mb-4 flex justify-between items-center">
              <p className="text-gray-600 dark:text-gray-400">
                {loading ? (
                  "Loading mentors..."
                ) : error ? (
                  <span className="text-red-500">{error}</span>
                ) : selectedDomain ? (
                  <>
                    Showing {mentors.length} {selectedSubfield ? selectedSubfield : selectedDomain} mentors
                  </>
                ) : (
                  "Showing no mentors"
                )}
              </p>
            </div>
            
            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {[1, 2, 3, 4].map((item) => (
                  <div key={item} className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 animate-pulse">
                    <div className="flex items-start space-x-3">
                      <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-full"></div>
                      <div className="flex-1">
                        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded mb-2 w-2/3"></div>
                        <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded mb-3 w-full"></div>
                        <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : error ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="text-red-500 mb-4 text-xl">⚠️</div>
                <p className="text-gray-500 dark:text-gray-400 mb-4">{error}</p>
                <button 
                  onClick={() => window.location.href = '/login'} 
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Go to Login
                </button>
              </div>
            ) : selectedDomain && mentors.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {mentors.map((mentor) => (
                  <MentorCard key={mentor.id} mentor={mentor} />
                ))}
              </div>
            ) : selectedDomain && mentors.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <p className="text-gray-500 dark:text-gray-400 mb-4">
                  No mentors found for {selectedSubfield || selectedDomain}
                </p>
                <button 
                  onClick={() => {
                    setSelectedDomain(null);
                    setSelectedSubfield(null);
                  }} 
                  className="text-blue-600 dark:text-blue-400 underline"
                >
                  Clear selection
                </button>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <p className="text-gray-500 dark:text-gray-400 mb-4">
                  Select a domain from the tabs above to view available mentors
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
      </div>
      <Footer />
    </div>
  );
};