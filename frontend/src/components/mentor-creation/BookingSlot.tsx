import React, { useState } from 'react';
import { Clock, Save, Plus, Trash2 } from 'lucide-react';
import axios from 'axios';
import { API_URL } from '../../App';
import { useNavigate } from 'react-router-dom';

interface BookingSlotData {
  days: string[];
  startTime: string;
  startPeriod: string;
  endTime: string;
  endPeriod: string;
}

export const BookingSlot = () => {
  const [bookingSlots, setBookingSlots] = useState<BookingSlotData[]>([{
    days: [],
    startTime: '',
    startPeriod: 'AM',
    endTime: '',
    endPeriod: 'AM'
  }]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const daysOfWeek = [
    'Monday', 'Tuesday', 'Wednesday', 'Thursday',
    'Friday', 'Saturday', 'Sunday'
  ];

  const handleDayToggle = (slotIndex: number, day: string) => {
    setBookingSlots(prev => {
      const newSlots = [...prev];
      const slot = newSlots[slotIndex];
      newSlots[slotIndex] = {
        ...slot,
        days: slot.days.includes(day)
          ? slot.days.filter(d => d !== day)
          : [...slot.days, day]
      };
      return newSlots;
    });
  };

  const handleTimeChange = (
    slotIndex: number,
    field: keyof BookingSlotData,
    value: string
  ) => {
    setBookingSlots(prev => {
      const newSlots = [...prev];
      newSlots[slotIndex] = {
        ...newSlots[slotIndex],
        [field]: value
      };
      return newSlots;
    });
  };

  const addNewSlot = () => {
    setBookingSlots(prev => [...prev, {
      days: [],
      startTime: '',
      startPeriod: 'AM',
      endTime: '',
      endPeriod: 'PM'
    }]);
  };

  const removeSlot = (slotIndex: number) => {
    setBookingSlots(prev => prev.filter((_, index) => index !== slotIndex));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    const token = localStorage.getItem('token');
    if (!token) {
      setError('Please log in to continue');
      setIsLoading(false);
      return;
    }

    // Validate all slots
    for (const slot of bookingSlots) {
      if (!slot.days || slot.days.length === 0) {
        setError('Please select at least one day for each slot');
        setIsLoading(false);
        return;
      }
      if (!slot.startTime || !slot.endTime) {
        setError('Please set start and end times for each slot');
        setIsLoading(false);
        return;
      }
    }

    try {
      // Send each slot as a separate POST request
      for (const slot of bookingSlots) {
        const payload = {
          days: slot.days,
          start_time: slot.startTime,
          start_period: slot.startPeriod,
          end_time: slot.endTime,
          end_period: slot.endPeriod
        };
        console.log(payload)
        const response = await axios.post(
          `${API_URL}/v1/bookingslots/create`,
          payload, // Send each slot individually as a JSON object
          {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          }
        );

        if (response.status !== 201) {
          throw new Error(`Failed to create booking slot: ${slot.startTime}-${slot.endTime}`);
        }
      }

      navigate('/explore');
    } catch (error: any) {
      setError(error.response?.data?.message || error.message || 'Error creating booking slots');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {bookingSlots.map((slot, slotIndex) => (
        <div key={slotIndex} className="space-y-6 border-b pb-6 last:border-b-0">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Booking Slot {slotIndex + 1}
            </h3>
            {bookingSlots.length > 1 && (
              <button
                type="button"
                onClick={() => removeSlot(slotIndex)}
                className="text-red-500 hover:text-red-700"
              >
                <Trash2 size={18} />
              </button>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Available Days
            </label>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
              {daysOfWeek.map((day) => (
                <button
                  key={day}
                  type="button"
                  onClick={() => handleDayToggle(slotIndex, day)}
                  className={`px-4 py-2 rounded-lg border text-sm ${
                    slot.days.includes(day)
                      ? 'bg-blue-600 text-white border-blue-600'
                      : 'bg-white text-gray-700 border-gray-300 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600'
                  }`}
                >
                  {day}
                </button>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Start Time
              </label>
              <div className="flex gap-2">
                <div className="relative flex-1">
                  <Clock size={20} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                  <input
                    type="text"
                    required
                    value={slot.startTime}
                    onChange={(e) => handleTimeChange(slotIndex, 'startTime', e.target.value)}
                    placeholder="HH:MM"
                    pattern="^(1[0-2]|0?[1-9]):[0-5][0-9]$"
                    className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <select
                  value={slot.startPeriod}
                  onChange={(e) => handleTimeChange(slotIndex, 'startPeriod', e.target.value)}
                  className="w-20 px-2 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
                >
                  <option value="AM">AM</option>
                  <option value="PM">PM</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                End Time
              </label>
              <div className="flex gap-2">
                <div className="relative flex-1">
                  <Clock size={20} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                  <input
                    type="text"
                    required
                    value={slot.endTime}
                    onChange={(e) => handleTimeChange(slotIndex, 'endTime', e.target.value)}
                    placeholder="HH:MM"
                    pattern="^(1[0-2]|0?[1-9]):[0-5][0-9]$"
                    className="w-full pl-10 pr-4 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <select
                  value={slot.endPeriod}
                  onChange={(e) => handleTimeChange(slotIndex, 'endPeriod', e.target.value)}
                  className="w-20 px-2 py-2 border rounded-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:ring-2 focus:ring-blue-500"
                >
                  <option value="AM">AM</option>
                  <option value="PM">PM</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      ))}

      <div>
        <button
          type="button"
          onClick={addNewSlot}
          className="flex items-center gap-2 text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
        >
          <Plus size={18} />
          Add Another Slot
        </button>
      </div>

      {error && <div className="text-red-500 text-sm">{error}</div>}

      <div className="flex justify-end">
        <button
          type="submit"
          disabled={isLoading}
          className={`bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2 ${
            isLoading ? 'opacity-50 cursor-not-allowed' : ''
          }`}
        >
          {isLoading ? 'Saving...' : 'Save Profile'}
          {!isLoading && <Save size={18} />}
        </button>
      </div>
    </form>
  );
};