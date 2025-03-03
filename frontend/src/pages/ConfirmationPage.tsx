import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { API_URL } from '../App';
import { Loader2, CheckCircle2, XCircle } from 'lucide-react';

export const ConfirmationPage = () => {
  const { token = '' } = useParams<{ token: string }>();
  const navigate = useNavigate();
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState<string>('');

  const handleConfirm = async () => {
    setStatus('loading');
    try {
      const response = await fetch(`${API_URL}/users/activate/${token}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        setStatus('success');
        setTimeout(() => navigate('/explore'), 2000); // Redirect after 2 seconds
      } else {
        const errorData = await response.json();
        setStatus('error');
        setErrorMessage(errorData.message || 'Failed to confirm token');
      }
    } catch (error) {
      setStatus('error');
      setErrorMessage('Something went wrong. Please try again later.');
    }
  };

  // Optional: Auto-confirm on page load
  useEffect(() => {
    if (token) {
      handleConfirm();
    }
  }, [token]);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center px-4">
      <div className="max-w-md w-full bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 text-center">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">
          Account Confirmation
        </h1>

        {status === 'idle' && (
          <div>
            <p className="text-gray-600 dark:text-gray-300 mb-6">
              Click the button below to confirm your account.
            </p>
            <button
              onClick={handleConfirm}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 transition duration-150"
            >
              Confirm Account
            </button>
          </div>
        )}

        {status === 'loading' && (
          <div className="flex flex-col items-center">
            <Loader2 className="w-12 h-12 text-blue-600 animate-spin mb-4" />
            <p className="text-gray-600 dark:text-gray-300">
              Confirming your account...
            </p>
          </div>
        )}

        {status === 'success' && (
          <div className="flex flex-col items-center">
            <CheckCircle2 className="w-16 h-16 text-green-500 mb-4" />
            <p className="text-gray-900 dark:text-white font-medium mb-2">
              Account Confirmed!
            </p>
            <p className="text-gray-600 dark:text-gray-300">
              Redirecting to Explore page in a moment...
            </p>
          </div>
        )}

        {status === 'error' && (
          <div className="flex flex-col items-center">
            <XCircle className="w-16 h-16 text-red-500 mb-4" />
            <p className="text-gray-900 dark:text-white font-medium mb-2">
              Confirmation Failed
            </p>
            <p className="text-gray-600 dark:text-gray-300 mb-4">
              {errorMessage}
            </p>
            <button
              onClick={handleConfirm}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 transition duration-150"
            >
              Try Again
            </button>
            <button
              onClick={() => navigate('/')}
              className="mt-2 text-blue-600 dark:text-blue-400 hover:underline"
            >
              Return to Home
            </button>
          </div>
        )}
      </div>
    </div>
  );
};