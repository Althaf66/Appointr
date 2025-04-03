import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { 
  ArrowLeft, 
  Calendar, 
  Clock, 
  Briefcase, 
  GraduationCap, 
  MapPin, 
  Star, 
  ChevronRight,
  Globe, 
  MessageCircle 
} from 'lucide-react';

interface MentorDetails {
  id: string;
  userid: string;
  name: string;
  country: string;
  language: string[];
  gigs: { title: string; description: string; expertise: string; discipline: string[] }[];
  education: { degree: string; field: string; institute: string; year_from: string; year_to: string }[];
  experience: { title: string; company: string; description: string; year_from: string; year_to: string }[];
  workingat: { title: string; company: string; totalyear: number };
  bookingslots: { days: string[]; start_time: string; end_time: string; start_period: string; end_period: string }[];
}

interface TimeSlot {
  time: string;
  period: string;
}

export const MentorDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [mentor, setMentor] = useState<MentorDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showBookingModal, setShowBookingModal] = useState(false);
  const [selectedDate, setSelectedDate] = useState('');
  const [selectedTimeSlot, setSelectedTimeSlot] = useState<TimeSlot | null>(null);

  const getAuthToken = () => localStorage.getItem('token');
  
  const getAxiosConfig = () => ({
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    }
  });

  const handleMessage = async () => {
    try {
      if (!mentor?.userid) {
        console.error('Mentor userId not available');
        return;
      }

      const response = await axios.post(
        'http://localhost:8080/v1/messages/conversations',
        {
          other_user_id: mentor.userid
        },
        getAxiosConfig()
      );
      
      // On successful response, navigate to messages page
      navigate('/messages');
    } catch (err) {
      console.error('Failed to create conversation:', err);
      // Optionally, you could show an error message to the user here
    }
  };

  useEffect(() => {
    const fetchMentorDetails = async () => {
      try {
        setLoading(true);
        const response = await axios.get(
          `http://localhost:8080/v1/mentors/u/${id}`,
          getAxiosConfig()
        );
        console.log('API response:', response.data.data);
        setMentor(response.data.data);
        
        // Set default selected date to first available day
        if (response.data.data.bookingslots?.length > 0) {
          const today = new Date();
          const dayNames = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'];
          const availableDays = response.data.data.bookingslots.flatMap((slot: { days: any; }) => slot.days);
          
          // Find the next available date based on bookingslots.days
          for (let i = 0; i < 7; i++) {
            const nextDate = new Date(today);
            nextDate.setDate(today.getDate() + i);
            const dayName = dayNames[nextDate.getDay()];
            
            if (availableDays.includes(dayName)) {
              setSelectedDate(formatDate(nextDate));
              break;
            }
          }
        }
        
        setLoading(false);
      } catch (err) {
        setError('Failed to load mentor details');
        setLoading(false);
      }
    };

    fetchMentorDetails();
  }, [id]);

  const formatDate = (date: Date): string => {
    const day = String(date.getDate()).padStart(2, '0');
    const month = date.toLocaleString('default', { month: 'short' });
    const year = date.getFullYear();
    return `${day} ${month} ${year}`;
  };

  const getAvailableDates = () => {
    if (!mentor?.bookingslots?.length) return [];
    
    const availableDays = mentor.bookingslots.flatMap(slot => slot.days);
    const today = new Date();
    const dates = [];
    
    for (let i = 0; i < 7; i++) {
      const nextDate = new Date(today);
      nextDate.setDate(today.getDate() + i);
      const dayName = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'][nextDate.getDay()];
      
      if (availableDays.includes(dayName)) {
        dates.push({
          date: nextDate,
          dayName: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'][nextDate.getDay()],
          day: nextDate.getDate(),
          month: nextDate.toLocaleString('default', { month: 'short' }),
          formatted: formatDate(nextDate)
        });
      }
    }
    
    return dates;
  };

  const getTimeSlots = () => {
    if (!mentor?.bookingslots?.length) return [];
    
    // Get day of the week for selected date
    const [day, month, year] = selectedDate.split(' ');
    const selectedDateObj = new Date(`${month} ${day} ${year}`);
    const dayName = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'][selectedDateObj.getDay()];
    
    // Find booking slots for this day
    const availableSlots = mentor.bookingslots.filter(slot => 
      slot.days.includes(dayName)
    );
    
    if (availableSlots.length === 0) return [];
    
    // Generate 30 min time slots between start and end time
    const timeSlots: TimeSlot[] = [];
    
    availableSlots.forEach(slot => {
      let startHour = parseInt(slot.start_time);
      let endHour = parseInt(slot.end_time);
      
      // Adjust for 12-hour format
      if (slot.start_period === 'PM' && startHour !== 12) startHour += 12;
      if (slot.end_period === 'PM' && endHour !== 12) endHour += 12;
      if (slot.start_period === 'AM' && startHour === 12) startHour = 0;
      if (slot.end_period === 'AM' && endHour === 12) endHour = 0;
      
      for (let hour = startHour; hour < endHour; hour++) {
        for (let minute = 0; minute < 60; minute += 30) {
          if (hour === endHour && minute > 0) continue;
          
          let displayHour = hour % 12;
          if (displayHour === 0) displayHour = 12;
          
          const period = hour >= 12 ? 'PM' : 'AM';
          const timeString = `${displayHour}:${minute === 0 ? '00' : minute}`;
          
          timeSlots.push({ time: timeString, period });
        }
      }
    });
    
    return timeSlots;
  };

  const handleBookSession = () => {
    if (!selectedTimeSlot) return;
    
    // Here you would implement the actual booking logic
    alert(`Booking confirmed for ${selectedDate} at ${selectedTimeSlot.time} ${selectedTimeSlot.period}`);
    setShowBookingModal(false);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
        <div className="animate-pulse max-w-4xl mx-auto bg-white dark:bg-gray-800 rounded-lg p-6">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-4"></div>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-2/3 mb-6"></div>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"></div>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"></div>
        </div>
      </div>
    );
  }

  if (error || !mentor) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 p-8 text-center">
        <p className="text-red-500">{error || 'Mentor not found'}</p>
        <button onClick={() => navigate('/explore')} className="mt-4 text-blue-600 underline">
          Back to Explore
        </button>
      </div>
    );
  }

  const availableDates = getAvailableDates();
  const timeSlots = getTimeSlots();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      {/* Main container with increased width */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <button 
          onClick={() => navigate('/explore')} 
          className="flex items-center text-blue-600 dark:text-blue-400 mb-6"
        >
          <ArrowLeft size={18} className="mr-2" /> Back to Mentors
        </button>
  
        <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
          {/* Profile Header */}
          <div className="border-b border-gray-200 dark:border-gray-700 p-6">
            <div className="flex items-center space-x-4">
              <img
                src={`https://api.dicebear.com/7.x/initials/svg?seed=${mentor.name}`}
                alt={mentor.name}
                className="w-20 h-20 rounded-full border-2 border-gray-200 dark:border-gray-700"
              />
              <div className="flex-1">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{mentor.name}</h1>
                
                <div className="flex items-center mt-1 text-gray-600 dark:text-gray-400">
                  {mentor.workingat && (
                    <div className="flex items-center mr-3">
                      <Briefcase size={16} className="mr-1" />
                      <span>{mentor.workingat.title} at {mentor.workingat.company}</span>
                    </div>
                  )}
                </div>
                
                <div className="flex items-center mt-1 text-gray-600 dark:text-gray-400">
                  <div className="flex items-center mr-3">
                    <MapPin size={16} className="mr-1" />
                    <span>{mentor.country}</span>
                  </div>
                  
                  <div className="flex items-center">
                    <Globe size={16} className="mr-1" />
                    <span>{mentor.language.join(', ')}</span>
                  </div>
                </div>
              </div>
  
              <div className="flex gap-3">
                <button 
                  onClick={() => setShowBookingModal(true)}
                  className="bg-teal-600 hover:bg-teal-700 text-white rounded-lg px-4 py-2 text-sm font-medium transition-colors"
                >
                  Book a Session
                </button>
                <button 
            onClick={handleMessage}
            className="bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg px-4 py-2 text-sm font-medium transition-colors"
          >
            Message
          </button>
              </div>
            </div>
          </div>
  
          {/* Main Content with adjusted grid */}
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Main content - 3 columns */}
            <div className="lg:col-span-3 p-6">
              {/* Expertise Tags */}
              <div className="mb-6">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Expertise</h2>
                <div className="flex flex-wrap gap-2">
                  {mentor.gigs?.flatMap(gig => [gig.expertise, ...gig.discipline])
                    .filter((value, index, self) => self.indexOf(value) === index)
                    .map((tag, index) => (
                      <span 
                        key={index} 
                        className="bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-3 py-1 text-sm rounded-full"
                      >
                        {tag}
                      </span>
                    ))}
                </div>
              </div>
              
              {/* Gigs */}
              {mentor.gigs?.length > 0 && (
                <section className="mb-6">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Services Offered</h2>
                  <div className="space-y-4">
                    {mentor.gigs.map((gig, index) => (
                      <div 
                        key={index} 
                        className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600"
                      >
                        <h3 className="font-medium text-gray-900 dark:text-white text-lg">{gig.title}</h3>
                        <p className="text-gray-600 dark:text-gray-400 mt-2">{gig.description}</p>
                      </div>
                    ))}
                  </div>
                </section>
              )}
  
              {/* Experience */}
              {(mentor.experience?.length > 0 || mentor.workingat) && (
                <section className="mb-6">
                  <div className="flex items-center mb-3">
                    <Briefcase size={18} className="mr-2 text-gray-600 dark:text-gray-400" />
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Professional Experience</h2>
                  </div>
                  
                  {mentor.workingat && (
                    <div className="mb-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600">
                      <div className="flex justify-between">
                        <div>
                          <h3 className="font-medium text-gray-900 dark:text-white">{mentor.workingat.title}</h3>
                          <p className="text-gray-600 dark:text-gray-400">{mentor.workingat.company}</p>
                        </div>
                        <div className="text-right">
                          <span className="inline-flex items-center rounded-full bg-green-100 dark:bg-green-900/30 px-3 py-1 text-xs font-medium text-green-800 dark:text-green-300">
                            Current
                          </span>
                          <p className="text-gray-500 dark:text-gray-400 text-sm mt-1">{mentor.workingat.totalyear} years</p>
                        </div>
                      </div>
                    </div>
                  )}
                  
                  {mentor.experience?.length > 0 && mentor.experience.map((exp, index) => (
                    <div 
                      key={index} 
                      className="mb-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600"
                    >
                      <div className="flex justify-between">
                        <div>
                          <h3 className="font-medium text-gray-900 dark:text-white">{exp.title}</h3>
                          <p className="text-gray-600 dark:text-gray-400">{exp.company}</p>
                        </div>
                        <p className="text-gray-500 dark:text-gray-400 text-sm">
                          {exp.year_from} - {exp.year_to}
                        </p>
                      </div>
                      <p className="text-gray-600 dark:text-gray-400 text-sm mt-2">{exp.description}</p>
                    </div>
                  ))}
                </section>
              )}
  
              {/* Education */}
              {mentor.education?.length > 0 && (
                <section className="mb-6">
                  <div className="flex items-center mb-3">
                    <GraduationCap size={18} className="mr-2 text-gray-600 dark:text-gray-400" />
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Education</h2>
                  </div>
                  
                  <div className="space-y-4">
                    {mentor.education.map((edu, index) => (
                      <div 
                        key={index} 
                        className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600"
                      >
                        <h3 className="font-medium text-gray-900 dark:text-white">{edu.degree} in {edu.field}</h3>
                        <p className="text-gray-600 dark:text-gray-400">{edu.institute}</p>
                        <p className="text-gray-500 dark:text-gray-400 text-sm">{edu.year_from} - {edu.year_to}</p>
                      </div>
                    ))}
                  </div>
                </section>
              )}
            </div>
  
            {/* Sidebar with booking slots - rightmost column */}
            <div className="p-6 lg:col-span-1">
              {mentor.bookingslots?.length > 0 && (
                <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-4 mb-6 border border-gray-200 dark:border-gray-600 sticky top-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center">
                      <Calendar size={18} className="mr-2 text-gray-600 dark:text-gray-400" />
                      <h3 className="font-medium text-gray-900 dark:text-white">Availability</h3>
                    </div>
                    <button 
                      onClick={() => setShowBookingModal(true)}
                      className="text-blue-600 dark:text-blue-400 text-sm hover:underline flex items-center"
                    >
                      View all <ChevronRight size={16} />
                    </button>
                  </div>
                  
                  <div className="space-y-3">
                    {mentor.bookingslots.slice(0, 3).map((slot, index) => (
                      <div 
                        key={index} 
                        className="flex items-center justify-between p-3 bg-white dark:bg-gray-800 rounded-lg shadow-sm"
                      >
                        <div>
                          <p className="text-gray-900 dark:text-white font-medium">{slot.days.join(', ')}</p>
                          <p className="text-gray-500 dark:text-gray-400 text-sm">
                            {slot.start_time} {slot.start_period} - {slot.end_time} {slot.end_period}
                          </p>
                        </div>
                        <button 
                          onClick={() => setShowBookingModal(true)}
                          className="text-teal-600 hover:text-teal-700 dark:text-teal-400 dark:hover:text-teal-300"
                        >
                          <Calendar size={18} />
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
  
        {/* Booking Modal */}
        {showBookingModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
            <div className="bg-white dark:bg-gray-800 rounded-xl max-w-md w-full mx-4 overflow-hidden">
              <div className="p-6">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">Book a Session</h3>
                  <button 
                    onClick={() => setShowBookingModal(false)} 
                    className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
                
                <div className="mb-6">
                  <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Available sessions</h4>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-4">Book 1:1 sessions from the options based on your needs</p>
                  
                  {/* Date Selection */}
                  <div className="flex overflow-x-auto space-x-2 py-2 mb-4">
                    {availableDates.map((date, index) => (
                      <button
                        key={index}
                        onClick={() => setSelectedDate(date.formatted)}
                        className={`flex-shrink-0 p-3 rounded-lg border ${
                          selectedDate === date.formatted 
                            ? 'border-teal-600 bg-teal-50 dark:bg-teal-900/20 dark:border-teal-500' 
                            : 'border-gray-200 dark:border-gray-700'
                        }`}
                      >
                        <p className="text-xs text-center text-gray-500 dark:text-gray-400">{date.dayName}</p>
                        <p className={`text-lg font-medium text-center ${
                          selectedDate === date.formatted 
                            ? 'text-teal-600 dark:text-teal-400' 
                            : 'text-gray-900 dark:text-white'
                        }`}>{date.day}</p>
                        <p className="text-xs text-center text-gray-500 dark:text-gray-400">{date.month}</p>
                        <p className="text-xs text-center text-green-500 mt-1">23 slots</p>
                      </button>
                    ))}
                  </div>
                  
                  {/* Time Slots */}
                  <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 flex justify-between items-center">
                    Available time slots
                    <ChevronRight size={16} className="text-gray-400" />
                  </h4>
                  
                  <div className="grid grid-cols-3 gap-3 mb-4">
                    {timeSlots.map((slot, index) => (
                      <button
                        key={index}
                        onClick={() => setSelectedTimeSlot(slot)}
                        className={`p-3 rounded-lg border text-center ${
                          selectedTimeSlot?.time === slot.time && selectedTimeSlot?.period === slot.period
                            ? 'border-teal-600 bg-teal-50 dark:bg-teal-900/20 dark:border-teal-500 text-teal-600 dark:text-teal-400' 
                            : 'border-gray-200 dark:border-gray-700 text-gray-900 dark:text-white'
                        }`}
                      >
                        {slot.time} {slot.period}
                      </button>
                    ))}
                  </div>
                  
                  <button
                    onClick={handleBookSession}
                    disabled={!selectedTimeSlot}
                    className={`w-full py-3 rounded-lg text-white font-medium ${
                      selectedTimeSlot 
                        ? 'bg-teal-600 hover:bg-teal-700' 
                        : 'bg-gray-300 dark:bg-gray-600 cursor-not-allowed'
                    }`}
                  >
                    Book Session for {selectedDate.split(' ').slice(0, 2).join(' ')} {selectedDate.includes('2025') ? '2025' : ''}
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}