import React from 'react';
import { MentorHeader } from '../components/mentor/MentorHeader';
import { MentorProfile } from '../components/mentor/MentorProfile';
import { BookingSection } from '../components/mentor/BookingSection';
import { Background } from '../components/mentor/Background';
import { SimilarMentors } from '../components/mentor/SimilarMentors';
import { Footer } from '../components/Footer';

export const MentorProfilePage = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <MentorHeader />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <MentorProfile />
            <Background />
            <SimilarMentors />
          </div>
          
          <div className="lg:col-span-1">
            <div className="sticky top-24">
              <BookingSection />
            </div>
          </div>
        </div>
      </main>
      
      <Footer />
    </div>
  );
};