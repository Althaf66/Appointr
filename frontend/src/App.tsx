import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
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
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/explore" element={<ExplorePage />} />
          <Route path="/mentor/:id" element={<MentorProfilePage />} />
          <Route path="/messages" element={<MessagesPage />} />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;