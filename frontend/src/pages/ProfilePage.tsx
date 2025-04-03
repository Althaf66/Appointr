import React from 'react';
import { MentorHeader } from '../components/mentor/MentorHeader';
import { ProfileHeader } from '../components/profile/ProfileHeader';
import { ProfileStats } from '../components/profile/ProfileStats';
import { UpcomingSessions } from '../components/profile/UpcomingSessions';
import { ProfileSettings } from '../components/profile/ProfileSettings';
import { Footer } from '../components/Footer';

export const ProfilePage = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <MentorHeader />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-6">
            <ProfileHeader />
            {/* <ProfileStats /> */}
            {/* <UpcomingSessions /> */}
          </div>
          <div className="lg:col-span-1">
            <div className="sticky top-24">
              {/* <ProfileSettings /> */}
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
};