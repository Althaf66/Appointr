import React, { useEffect, useState } from 'react';
import { MentorHeader } from '../components/mentor/MentorHeader';
import { Footer } from '../components/Footer';
import axios from 'axios';

// Types remain unchanged
interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

interface Meeting {
  id: number;
  userid: number;
  mentorid: number;
  day: string;
  date: string;
  start_time: string;
  start_period: string;
  isconfirm: boolean;
  ispaid: boolean;
  iscompleted: boolean;
  amount: number;
  link: string;
  menteeName?: string;
  mentorName?: string;
}

export const ProfilePage = () => {
  const [user, setUser] = useState<User | null>(null);
  const [meetings, setMeetings] = useState<Meeting[]>([]);
  const [unpaidMeetings, setUnpaidMeetings] = useState<Meeting[]>([]);
  const [loading, setLoading] = useState(true);
  const [unpaidLoading, setUnpaidLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [unpaidError, setUnpaidError] = useState<string | null>(null);
  const [paidMeetings, setPaidMeetings] = useState<Meeting[]>([]);
  const [paidLoading, setPaidLoading] = useState(true);
  const [paidError, setPaidError] = useState<string | null>(null);
  const [mentorPaidMeetings, setMentorPaidMeetings] = useState<Meeting[]>([]);
  const [mentorPaidLoading, setMentorPaidLoading] = useState(true);
  const [mentorPaidError, setMentorPaidError] = useState<string | null>(null);

  const token = localStorage.getItem('token');

  function getUserIdFromJWT(token: string): number | null {
    try {
      const payload = token.split('.')[1];
      const decodedPayload = atob(payload);
      const jsonPayload = JSON.parse(decodedPayload);
      return jsonPayload.sub || null;
    } catch (error) {
      console.error('Error decoding JWT token:', error);
      return null;
    }
  }

  const userId = token ? getUserIdFromJWT(token) : null;

  const [message, setMessage] = useState("");

  useEffect(() => {
    const query = new URLSearchParams(window.location.search);
    if (query.get("success")) {
      setMessage("Booked! You will receive an email confirmation.");
    }
    if (query.get("canceled")) {
      setMessage("Booking canceled -- checkout when you're ready.");
    }
  }, []);

  useEffect(() => {
    axios.defaults.headers.common['Content-Type'] = 'application/json';
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    }
  }, [token]);

  useEffect(() => {
    const fetchPaidMeetings = async () => {
      try {
        setPaidLoading(true);
        if (!token || !userId) throw new Error('Authentication required');
  
        const config = { headers: { Authorization: `Bearer ${token}` } };
        const paidResponse = await axios.get(
          `http://localhost:8080/v1/meetings/user-not-completed/${userId}`,
          config
        );
        let paidData = Array.isArray(paidResponse.data.data)
          ? paidResponse.data.data
          : paidResponse.data.data ? [paidResponse.data.data] : [];
        const paidWithNames = await fetchUserNames(paidData);
        setPaidMeetings(paidWithNames);
  
      } catch (err: any) {
        setPaidError(err.message || 'An error occurred while fetching paid sessions');
        console.error('Error fetching paid meetings:', err);
      } finally {
        setPaidLoading(false);
      }
    };
    fetchPaidMeetings();
  }, [token, userId]);

  useEffect(() => {
    const fetchMentorPaidMeetings = async () => {
      try {
        setMentorPaidLoading(true);
        if (!token || !userId) throw new Error('Authentication required');
  
        const config = { headers: { Authorization: `Bearer ${token}` } };
        const mentorPaidResponse = await axios.get(
          `http://localhost:8080/v1/meetings/mentor-not-completed/${userId}`,
          config
        );
        let mentorPaidData = Array.isArray(mentorPaidResponse.data.data)
          ? mentorPaidResponse.data.data
          : mentorPaidResponse.data.data ? [mentorPaidResponse.data.data] : [];
        const mentorPaidWithNames = await fetchUserNames(mentorPaidData);
        setMentorPaidMeetings(mentorPaidWithNames);
  
      } catch (err: any) {
        setMentorPaidError(err.message || 'An error occurred while fetching mentor paid sessions');
        console.error('Error fetching mentor paid meetings:', err);
      } finally {
        setMentorPaidLoading(false);
      }
    };
    fetchMentorPaidMeetings();
  }, [token, userId]);

  const fetchUserNames = async (meetingsData: Meeting[]): Promise<Meeting[]> => {
    if (!token || !meetingsData.length) return meetingsData;

    const config = {
      headers: {
        Authorization: `Bearer ${token}`
      }
    };

    const meetingsWithNames = await Promise.all(
      meetingsData.map(async (meeting) => {
        const userIdToFetch = meeting.userid;
        const mentorIdToFetch = meeting.mentorid;
        const updatedMeeting = { ...meeting };

        try {
          const menteeResponse = await axios.get(
            `http://localhost:8080/v1/users/${userIdToFetch}`,
            config
          );
          updatedMeeting.menteeName = menteeResponse.data.data.username;

          if (mentorIdToFetch !== userId) {
            const mentorResponse = await axios.get(
              `http://localhost:8080/v1/users/${mentorIdToFetch}`,
              config
            );
            updatedMeeting.mentorName = mentorResponse.data.data.username;
          }

          return updatedMeeting;
        } catch (err) {
          console.error(`Error fetching user data:`, err);
          return {
            ...updatedMeeting,
            menteeName: updatedMeeting.menteeName || `User #${userIdToFetch}`,
            mentorName: updatedMeeting.mentorName || `Mentor #${mentorIdToFetch}`
          };
        }
      })
    );

    return meetingsWithNames;
  };

  const handlePayment = async (meetingId: number) => {
    try {
      if (!token) {
        throw new Error('No authentication token found');
      }

      const meeting = unpaidMeetings.find(m => m.id === meetingId);
      if (!meeting || !user) {
        throw new Error('Meeting or user data not found');
      }

      const config = {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        }
      };

      const paymentData = {
        price: meeting.amount,
        currency: 'inr',
        email: user.email,
        menteeName: meeting.menteeName,
        mentorName: meeting.mentorName,
        id: meeting.id,
      };

      const response = await axios.post(
        `http://localhost:8080/v1/payment/create-checkout-session`,
        paymentData,
        config
      );

      if (response.data.url) {
        window.location.href = response.data.url; // Redirect to Stripe checkout
      }

    } catch (err: any) {
      if (err.response && err.response.status === 0) {
        setUnpaidError('CORS error: The server is not allowing cross-origin requests. Please contact your administrator.');
      } else {
        setUnpaidError(err.message || 'Failed to process payment');
      }
      console.error('Error processing payment:', err);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        if (!token || !userId) throw new Error('Authentication required');
        
        const config = { headers: { Authorization: `Bearer ${token}` } };
        const userResponse = await axios.get(`http://localhost:8080/v1/users/${userId}`, config);
        setUser(userResponse.data.data);

        const meetingsResponse = await axios.get(
          `http://localhost:8080/v1/meetings/mentor-not-confirm/${userId}`,
          config
        );
        let meetingsData = Array.isArray(meetingsResponse.data.data)
          ? meetingsResponse.data.data
          : meetingsResponse.data.data ? [meetingsResponse.data.data] : [];
        const meetingsWithNames = await fetchUserNames(meetingsData);
        setMeetings(meetingsWithNames);

      } catch (err: any) {
        setError(err.message || 'An error occurred while fetching data');
        console.error('Error fetching data:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [token, userId]);

  useEffect(() => {
    const fetchUnpaidMeetings = async () => {
      try {
        setUnpaidLoading(true);
        if (!token || !userId) throw new Error('Authentication required');

        const config = { headers: { Authorization: `Bearer ${token}` } };
        const unpaidResponse = await axios.get(
          `http://localhost:8080/v1/meetings/user-not-paid/${userId}`,
          config
        );
        let unpaidData = Array.isArray(unpaidResponse.data.data)
          ? unpaidResponse.data.data
          : unpaidResponse.data.data ? [unpaidResponse.data.data] : [];
        const unpaidWithNames = await fetchUserNames(unpaidData);
        setUnpaidMeetings(unpaidWithNames);

      } catch (err: any) {
        setUnpaidError(err.message || 'An error occurred while fetching unpaid sessions');
        console.error('Error fetching unpaid meetings:', err);
      } finally {
        setUnpaidLoading(false);
      }
    };
    fetchUnpaidMeetings();
  }, [token, userId]);

  const handleConfirmMeeting = async (meetingId: number) => {
    try {
      if (!token) throw new Error('No authentication token found');
      const config = { headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' } };
      await axios.put(`http://localhost:8080/v1/meetings/confirm/${meetingId}`, {}, config);
      setMeetings(meetings.filter(meeting => meeting.id !== meetingId));
    } catch (err: any) {
      setError(err.message || 'Failed to confirm meeting');
      console.error('Error confirming meeting:', err);
    }
  };

  const formatDate = (dateString: string) => dateString;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <MentorHeader />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-6">
          {user && (
            <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">My Profile</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400"><span className="font-medium">Username:</span> {user.username}</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400"><span className="font-medium">Email:</span> {user.email}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-400"><span className="font-medium">Member Since:</span> {new Date(user.created_at).toLocaleDateString()}</p>
                </div>
              </div>
            </div>
          )}
  
          {/* Sessions Ready for Payment - Only show if there are meetings after loading */}
          {!unpaidLoading && unpaidMeetings && unpaidMeetings.length > 0 && (
            <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
              <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">Sessions Ready for Payment</h3>
              </div>
              {unpaidLoading ? (
                <div className="p-6 text-center"><p className="text-gray-500 dark:text-gray-400">Loading payment-ready sessions...</p></div>
              ) : unpaidError ? (
                <div className="p-6 text-center"><p className="text-red-500">{unpaidError}</p></div>
              ) : (
                <div className="divide-y divide-gray-200 dark:divide-gray-700">
                  {unpaidMeetings.map((meeting) => (
                    <div key={meeting.id} className="p-6">
                      <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                        <div>
                          <h4 className="font-medium text-gray-900 dark:text-white">{meeting.day} - {formatDate(meeting.date)}</h4>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Time: {meeting.start_time} {meeting.start_period}</p>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Mentor: {meeting.mentorName || `Mentor #${meeting.mentorid}`}</p>
                          <div className="mt-2">
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-green-100 text-green-800">Mentor Confirmed</span>
                          </div>
                        </div>
                        <div className="flex flex-col items-end">
                          <p className="mb-2 text-lg font-semibold text-gray-900 dark:text-white">₹{meeting.amount}</p>
                          <button
                            onClick={() => handlePayment(meeting.id)}
                            className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
                          >
                            Pay Now
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
  
          {/* Paid Sessions - Only show if there are meetings after loading */}
          {!paidLoading && paidMeetings && paidMeetings.length > 0 && (
            <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
              <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">Paid Sessions</h3>
              </div>
              {paidLoading ? (
                <div className="p-6 text-center"><p className="text-gray-500 dark:text-gray-400">Loading paid sessions...</p></div>
              ) : paidError ? (
                <div className="p-6 text-center"><p className="text-red-500">{paidError}</p></div>
              ) : (
                <div className="divide-y divide-gray-200 dark:divide-gray-700">
                  {paidMeetings.map((meeting) => (
                    <div key={meeting.id} className="p-6">
                      <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                        <div>
                          <h4 className="font-medium text-gray-900 dark:text-white">{meeting.day} - {formatDate(meeting.date)}</h4>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Time: {meeting.start_time} {meeting.start_period}</p>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Mentor: {meeting.mentorName || `Mentor #${meeting.mentorid}`}</p>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Amount: ₹{meeting.amount}</p>
                          {meeting.link && (
                            <a 
                              href={meeting.link} 
                              target="_blank" 
                              rel="noopener noreferrer"
                              className="mt-1 text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-600"
                            >
                              Meeting Link
                            </a>
                          )}
                        </div>
                        <div className="flex items-center">
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-green-100 text-green-800">
                            Paid
                          </span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

        {!mentorPaidLoading && mentorPaidMeetings && mentorPaidMeetings.length > 0 && (
          <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
            <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white">Upcoming Sessions</h3>
            </div>
            {mentorPaidLoading ? (
              <div className="p-6 text-center"><p className="text-gray-500 dark:text-gray-400">Loading mentor paid sessions...</p></div>
            ) : mentorPaidError ? (
              <div className="p-6 text-center"><p className="text-red-500">{mentorPaidError}</p></div>
            ) : (
              <div className="divide-y divide-gray-200 dark:divide-gray-700">
                {mentorPaidMeetings.map((meeting) => (
                  <div key={meeting.id} className="p-6">
                    <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                      <div>
                        <h4 className="font-medium text-gray-900 dark:text-white">{meeting.day} - {formatDate(meeting.date)}</h4>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Time: {meeting.start_time} {meeting.start_period}</p>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Mentee: {meeting.menteeName || `User #${meeting.userid}`}</p>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Amount: ₹{meeting.amount}</p>
                        {meeting.link && (
                          <a 
                            href={meeting.link} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="mt-1 text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-600"
                          >
                            Meeting Link
                          </a>
                        )}
                      </div>
                      <div className="flex items-center">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-green-100 text-green-800">
                          Paid
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
  
          {/* Upcoming Sessions - Only show if there are meetings after loading */}
          {!loading && meetings && meetings.length > 0 && (
            <div className="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
              <div className="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">Need Confirmation</h3>
              </div>
              {loading ? (
                <div className="p-6 text-center"><p className="text-gray-500 dark:text-gray-400">Loading sessions...</p></div>
              ) : error ? (
                <div className="p-6 text-center"><p className="text-red-500">{error}</p></div>
              ) : (
                <div className="divide-y divide-gray-200 dark:divide-gray-700">
                  {meetings.map((meeting) => (
                    <div key={meeting.id} className="p-6">
                      <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                        <div>
                          <h4 className="font-medium text-gray-900 dark:text-white">{meeting.day} - {formatDate(meeting.date)}</h4>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Time: {meeting.start_time} {meeting.start_period}</p>
                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Client Name: {meeting.menteeName || `User #${meeting.userid}`}</p>
                        </div>
                        <div className="flex justify-end">
                          <div className="text-center">
                            <p className="mb-2 text-sm text-gray-500 dark:text-gray-400">Are you ready to confirm?</p>
                            <button
                              onClick={() => handleConfirmMeeting(meeting.id)}
                              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                              Confirm Session
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </main>
      <Footer />
    </div>
  );
};