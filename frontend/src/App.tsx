import React from 'react';
import { ProfilePage } from './pages/ProfilePage';
import { LandingPage } from './pages/LandingPage';
import { LoginPage } from './pages/auth/LoginPage';
import { SignupPage } from './pages/auth/SignupPage';
import { ThemeProvider } from './context/ThemeContext';
import { ExplorePage } from './pages/ExplorePage';
import { MentorProfilePage } from './pages/MentorProfilePage';
import { MessagesPage } from './pages/MessagesPage';

function App() {
  return (
    <ThemeProvider>
      <LandingPage />
      <LoginPage />
      <SignupPage />
      <ExplorePage />
      <MentorProfilePage />
      <MessagesPage />
      <ProfilePage />
    </ThemeProvider>
  );
}

export default App;