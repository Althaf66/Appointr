import React from 'react';
import { MapPin, Briefcase } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

// Adjusted Mentor interface to use id (assuming this is what’s available)
interface Mentor {
  id: string; // Changed to string to match typical userid format; adjust if it’s a number
  userid: string;
  name: string;
  avatar?: string;
  country?: string;
  hourlyRate?: number;
  gigs?: { title: string; description: string; discipline: string[] }[];
  workingat?: { company: string; totalyear: number };
}

interface MentorCardProps {
  mentor: Mentor;
}

export const MentorCard = ({ mentor }: MentorCardProps) => {
  const navigate = useNavigate();

  console.log('Mentor data in MentorCard:', mentor); // Debug log

  // Format hourly rate with currency symbol
  const formattedRate = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 0,
  }).format(mentor.hourlyRate || 0);

  // Use mentor.id for navigation
  const handleCardClick = () => {
    navigate(`/mentor/${mentor.userid}`); // Use id instead of userid
  };

  const primaryGig = mentor.gigs && mentor.gigs.length > 0 ? mentor.gigs[0] : null;

  return (
    <div 
      className="bg-white dark:bg-gray-800 rounded-lg shadow-sm hover:shadow-md transition-all duration-300 p-4 border border-gray-200 dark:border-gray-700 cursor-pointer flex flex-col"
      onClick={handleCardClick}
    >
      <div className="flex items-start space-x-3">
        <img
          src={mentor.avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${mentor.name}`}
          alt={mentor.name}
          className="w-12 h-12 rounded-full object-cover border-2 border-gray-200 dark:border-gray-700 flex-shrink-0"
        />
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-base text-gray-900 dark:text-white truncate">{mentor.name}</h3>
          <p className="text-sm text-gray-600 dark:text-gray-400 truncate">
            {primaryGig ? primaryGig.title : 'No title available'}
          </p>
        </div>
      </div>
      <p className="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
        {primaryGig ? primaryGig.description : 'No description available'}
      </p>
      <div className="mt-2 space-y-1 text-sm text-gray-500 dark:text-gray-400 flex-1">
        {mentor.country && (
          <div className="flex items-center">
            <MapPin size={14} className="mr-1 flex-shrink-0" />
            <span className="truncate">{mentor.country}</span>
          </div>
        )}
        {mentor.workingat && (
          <div className="flex items-center">
            <Briefcase size={14} className="mr-1 flex-shrink-0" />
            <span className="truncate">
              {mentor.workingat.company} - {mentor.workingat.totalyear} years
            </span>
          </div>
        )}
      </div>
      {primaryGig && primaryGig.discipline && primaryGig.discipline.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {primaryGig.discipline.map((discipline, index) => (
            <span
              key={index}
              className="bg-purple-50 dark:bg-purple-900/ W30 text-purple-700 dark:text-purple-300 text-xs px-2 py-1 rounded-full"
            >
              {discipline}
            </span>
          ))}
        </div>
      )}
      <div 
        className="mt-3 pt-3 border-t dark:border-gray-700 flex items-center justify-between"
        onClick={(e) => e.stopPropagation()}
      >
        <span className="text-sm font-medium text-gray-900 dark:text-white">
          {mentor.hourlyRate ? `${formattedRate}/hr` : 'Rate TBD'}
        </span>
        <Link
          to={`/message/${mentor.id}`} // Use id instead of userid
          className="text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/30 px-3 py-1 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-800/30 transition-colors text-sm"
        >
          Message
        </Link>
      </div>
    </div>
  );
};