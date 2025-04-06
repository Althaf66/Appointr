import React from 'react';
import { MapPin, Briefcase } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

interface Mentor {
  id: string;
  userid: string;
  name: string;
  avatar?: string;
  country?: string;
  gigs?: { title: string; description: string; amount: number; discipline: string[] }[];
  workingat?: { company: string; totalyear: number };
}

interface MentorCardProps {
  mentor: Mentor;
}

export const MentorCard = ({ mentor }: MentorCardProps) => {
  const navigate = useNavigate();

  console.log('Mentor data in MentorCard:', mentor); // Debug log

  // Format hourly rate with currency symbol
  const primaryGig = mentor.gigs && mentor.gigs.length > 0 ? mentor.gigs[0] : null;
  const formattedRate = primaryGig?.amount
    ? new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'INR',
        minimumFractionDigits: 0,
      }).format(primaryGig.amount)
    : 'Rate TBD';

  // Use mentor.id for navigation
  const handleCardClick = () => {
    navigate(`/mentor/${mentor.userid}`);
  };

  return (
    <div
      className="bg-white dark:bg-gray-800 rounded-lg shadow-sm hover:shadow-lg transition-all duration-300 p-5 border border-gray-200 dark:border-gray-700 cursor-pointer flex flex-col"
      onClick={handleCardClick}
    >
      <div className="flex items-start justify-between space-x-4">
        <div className="flex items-start space-x-4 flex-1">
          <img
            src={mentor.avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${mentor.name}`}
            alt={mentor.name}
            className="w-14 h-14 rounded-full object-cover border-2 border-gray-200 dark:border-gray-700 flex-shrink-0"
          />
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-lg text-gray-900 dark:text-white truncate">{mentor.name}</h3>
            <p className="text-sm font-medium text-gray-700 dark:text-gray-300 truncate">
              {primaryGig ? primaryGig.title : 'No title available'}
            </p>
          </div>
        </div>
        {mentor.workingat && (
          <div className="text-sm text-gray-600 dark:text-gray-400 text-right flex-shrink-0">
            <div className="flex items-center justify-end">
              <Briefcase size={14} className="mr-1 flex-shrink-0" />
              <span className="truncate">
                {mentor.workingat.company} - {mentor.workingat.totalyear} yrs
              </span>
            </div>
          </div>
        )}
      </div>

      <p className="mt-3 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
        {primaryGig ? primaryGig.description : 'No description available'}
      </p>

      <div className="mt-3 space-y-2 text-sm text-gray-500 dark:text-gray-400">
        {mentor.country && (
          <div className="flex items-center">
            <MapPin size={14} className="mr-1 flex-shrink-0" />
            <span className="truncate">{mentor.country}</span>
          </div>
        )}
      </div>

      {primaryGig && primaryGig.discipline && primaryGig.discipline.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1.5">
          {primaryGig.discipline.map((discipline, index) => (
            <span
              key={index}
              className="bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs font-medium px-2.5 py-1 rounded-full"
            >
              {discipline}
            </span>
          ))}
        </div>
      )}

      <div className="mt-4 pt-4 border-t dark:border-gray-700">
        <span className="text-base font-semibold text-gray-900 dark:text-white">
          {formattedRate} <span className="text-xs font-normal text-gray-500 dark:text-gray-400"></span>
        </span>
      </div>
    </div>
  );
};