import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Code2, ShieldPlus, Scale, Tractor, GraduationCap, Award, HeartHandshake, HeartPulse, Handshake, Gem } from 'lucide-react';
import { Mentor } from '../types';
import { API_URL } from '../App';

const iconMap = {
  Code2,
  ShieldPlus,
  Scale,
  Tractor,
  GraduationCap,
  Award,
  HeartHandshake,
  Gem,
  HeartPulse,
  Handshake
};

// interface Mentor {
//   id: string;
//   name: string;
//   avatar?: string;
//   role?: string;
//   company?: string;
//   location?: string;
//   rating?: number;
//   bio?: string;
//   expertise?: string[];
//   disciplines?: string[];
//   yearsOfExperience?: number;
//   hourlyRate?: number;
//   isAvailableNow?: boolean;
//   availability?: string;
//   availabilityDetails?: string;
//   education?: { degree: string; institution: string; year: string }[];
//   sessionTypes?: string[];
// }

interface DomainTabsProps {
  onMentorsUpdate: (mentors: Mentor[], domain: string | null, subfield: string | null) => void;
}

export const DomainTabs: React.FC<DomainTabsProps> = ({ onMentorsUpdate }) => {
  interface Domain {
    id: string;
    name: string;
    icon_svg: string;
  }

  interface Subfield {
    id: string;
    name: string;
    subfield: string;
  }

  const [domains, setDomains] = useState<Domain[]>([]);
  const [subfields, setSubfields] = useState<Subfield[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [selectedSubfield, setSelectedSubfield] = useState<string | null>(null);

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

  // Fetch domains on component mount
  useEffect(() => {
    const fetchDomains = async () => {
      try {
        const response = await axios.get(
          `${API_URL}/expertise`,
          getAxiosConfig()
        );
        setDomains(response.data.data || []);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching domains:', error);
        setDomains([]);
        setLoading(false);
      }
    };

    fetchDomains();
  }, []);

  // Fetch subfields when a domain is selected
  useEffect(() => {
    const fetchSubfields = async () => {
      if (!selectedDomain) {
        setSubfields(null);
        return;
      }

      try {
        setLoading(true);
        const response = await axios.get(
          `${API_URL}/discipline/${selectedDomain}`,
          getAxiosConfig()
        );
        setSubfields(response.data.data || []);
        setLoading(false);
        
        // When a domain is selected, fetch mentors for that expertise
        fetchMentorsByExpertise(selectedDomain);
        
        // Reset selected subfield when domain changes
        setSelectedSubfield(null);
      } catch (error) {
        console.error('Error fetching subfields:', error);
        setSubfields([]);
        setLoading(false);
      }
    };

    fetchSubfields();
  }, [selectedDomain]);

  // Fetch mentors by expertise (domain)
  const fetchMentorsByExpertise = async (expertise: string) => {
    try {
      setLoading(true);
      const response = await axios.get(
        `${API_URL}/mentors/exp/${expertise}`,
        getAxiosConfig()
      );
      const mentors = response.data.data || [];
      setLoading(false);
      
      // Update parent component with mentors data
      onMentorsUpdate(mentors, expertise, null);
    } catch (error) {
      console.error('Error fetching mentors by expertise:', error);
      onMentorsUpdate([], expertise, null);
      setLoading(false);
    }
  };

  // Fetch mentors by discipline (subfield)
  const fetchMentorsByDiscipline = async (discipline: string) => {
    try {
      setLoading(true);
      const response = await axios.get(
        `${API_URL}/mentors/dis/${discipline}`,
        getAxiosConfig()
      );
      const mentors = response.data.data || [];
      setLoading(false);
      
      // Update parent component with mentors data
      onMentorsUpdate(mentors, selectedDomain, discipline);
    } catch (error) {
      console.error('Error fetching mentors by discipline:', error);
      onMentorsUpdate([], selectedDomain, discipline);
      setLoading(false);
    }
  };

  // Handle domain selection
  const handleDomainClick = (domainName: string) => {
    setSelectedDomain(domainName);
    setSelectedSubfield(null); // Reset subfield selection when domain changes
  };

  // Handle subfield selection
  const handleSubfieldClick = (subfield: string) => {
    setSelectedSubfield(subfield);
    fetchMentorsByDiscipline(subfield);
  };

  if (loading && !selectedDomain && !selectedSubfield) {
    return (
      <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8 py-4">
            <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-8 w-24 rounded"></div>
            <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-8 w-24 rounded"></div>
            <div className="animate-pulse bg-gray-200 dark:bg-gray-700 h-8 w-24 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Domain Tabs */}
        <div className="flex space-x-8 overflow-x-auto py-4">
          {domains.map((domain) => {
            const Icon = iconMap[domain.icon_svg as keyof typeof iconMap] || Code2;
            return (
              <button
                key={domain.id}
                onClick={() => handleDomainClick(domain.name)}
                className={`flex items-center space-x-2 px-1 py-2 border-b-2 ${
                  selectedDomain === domain.name
                    ? 'border-blue-600 text-gray-900 dark:text-white'
                    : 'border-transparent hover:border-blue-600 dark:hover:border-blue-400 text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
                } whitespace-nowrap`}
              >
                <Icon size={20} />
                <span>{domain.name}</span>
              </button>
            );
          })}
        </div>

        {/* Subfields Display */}
        {selectedDomain && !loading && subfields && subfields.length > 0 && (
          <div className="py-4">
            <div className="flex space-x-4 overflow-x-auto py-2">
              {subfields.map((subfield) => (
                <button
                  key={subfield.id}
                  onClick={() => handleSubfieldClick(subfield.subfield)}
                  className={`px-4 py-2 ${
                    selectedSubfield === subfield.subfield
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  } rounded-full whitespace-nowrap`}
                >
                  {subfield.subfield}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};