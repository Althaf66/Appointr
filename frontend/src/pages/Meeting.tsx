import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';

export const Meeting = () => {
  const [roomId, setRoomId] = useState('');
  const navigate = useNavigate();

  const createRoom = async () => {
    // Generate a random room ID if none provided
    const newRoomId = roomId || Math.random().toString(36).substring(2, 7);
    
    // Generate a random user ID
    const userId = 'user-' + Math.random().toString(36).substring(2, 7);
    
    try {
      const response = await fetch('http://localhost:8080/video/create-room', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ roomId: newRoomId, userId }),
      });
      
      if (response.ok) {
        const data = await response.json();
        // Store user ID in localStorage
        localStorage.setItem('userId', userId);
        // Navigate to the room
        navigate(`/room/${data.roomId}`);
      } else {
        alert('Failed to create room');
      }
    } catch (error) {
      console.error('Error creating room:', error);
      alert('Error connecting to server');
    }
  };

  const joinRoom = () => {
    if (!roomId) {
      alert('Please enter a room ID');
      return;
    }
    navigate(`/room/${roomId}`);
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <div className="p-8 bg-white rounded-lg shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-6 text-center">Video Meeting App</h1>
        
        <div className="mb-4">
          <label htmlFor="roomId" className="block text-sm font-medium text-gray-700 mb-1">
            Room ID (optional for create)
          </label>
          <input
            type="text"
            id="roomId"
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md"
            placeholder="Enter room ID"
          />
        </div>
        
        <div className="flex space-x-4">
          <button
            onClick={createRoom}
            className="flex-1 bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600"
          >
            Create Room
          </button>
          <button
            onClick={joinRoom}
            className="flex-1 bg-green-500 text-white py-2 px-4 rounded-md hover:bg-green-600"
          >
            Join Room
          </button>
        </div>
      </div>
    </div>
  );
};