import React from 'react';
import { Navbar } from '../components/Navbar';
import { HeroSection } from '../components/landing/HeroSection';
import { HowItWorks } from '../components/landing/HowItWorks';
import { PlatformStats } from '../components/landing/PlatformStats';
// import { Testimonials } from '../components/landing/Testimonials';
import { Footer } from '../components/Footer';
import { AboutStats } from '../components/landing/AboutStats';

export const LandingPage = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Navbar />
      <HeroSection />
      <HowItWorks />
      <AboutStats />
      <PlatformStats />
      <Footer />
    </div>
  );
};