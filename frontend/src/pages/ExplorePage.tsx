import React from 'react';
import { StickyHeader } from '../components/StickyHeader';
import { Navbar } from '../components/Navbar';
import { DomainTabs } from '../components/DomainTabs';
import { FilterSidebar } from '../components/explore/FilterSidebar';
import { ExploreHeader } from '../components/explore/ExploreHeader';
import { MentorCard } from '../components/MentorCard';
import { Footer } from '../components/Footer';
import { engineeringMentors } from '../data/engineeringMentors';

export const ExplorePage = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
      <StickyHeader />
      <Navbar />
      <DomainTabs />
      <ExploreHeader />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex gap-8">
          <FilterSidebar />
          
          <div className="flex-1">
            <div className="mb-4">
              <p className="text-gray-600 dark:text-gray-400">
                Showing {engineeringMentors.length} engineering mentors
              </p>
            </div>
            
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {engineeringMentors.map((mentor) => (
                <MentorCard key={mentor.id} mentor={mentor} />
              ))}
            </div>
          </div>
        </div>
      </div>
      
      <Footer />
    </div>
  );
};