import React, { useState } from 'react';
import { GigDetails } from '../components/mentor-creation/GigDetails';
import { ExperienceDetails } from '../components/mentor-creation/ExperienceDetails';
import { EducationDetails } from '../components/mentor-creation/EducationDetails';
import { WorkingNowDetails } from '../components/mentor-creation/WorkingNowDetails';
import { BookingSlot } from '../components/mentor-creation/BookingSlot';

const steps = [
  { id: 'gig', title: 'Gig Information' },
  { id: 'experience', title: 'Experience' },
  { id: 'education', title: 'Education' },
  { id: 'working', title: 'Current Work' },
  { id: 'bookingslot', title: 'Booking Slot' },
];

export const MentorCreationPage = () => {
  const [currentStep, setCurrentStep] = useState(0);

  const handleNext = () => {
    setCurrentStep((prev) => Math.min(prev + 1, steps.length - 1));
  };

  // const handlePrevious = () => {
  //   setCurrentStep((prev) => Math.max(prev - 1, 0));
  // };

  // case 0:
    // return <BasicDetails onNext={handleNext} />;
  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return <GigDetails onNext={handleNext} />;
      case 1:
        return <ExperienceDetails onNext={handleNext} />;
      case 2:
        return <EducationDetails onNext={handleNext} />;
      case 3:
        return <WorkingNowDetails onNext={handleNext} />;
      case 4:
        return <BookingSlot />;
      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-3xl mx-auto px-4 py-12">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
          {/* Progress bar */}
          <div className="px-8 py-4 border-b border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-4">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                Create Gig
              </h1>
              <span className="text-sm text-gray-500 dark:text-gray-400">
                Step {currentStep + 1} of {steps.length}
              </span>
            </div>
            <div className="relative">
              <div className="overflow-hidden h-2 rounded-full bg-gray-200 dark:bg-gray-700">
                <div
                  className="h-2 rounded-full bg-blue-600 transition-all duration-300"
                  style={{ width: `${((currentStep + 1) / steps.length) * 100}%` }}
                />
              </div>
              <div className="flex justify-between mt-2">
                {steps.map((step, index) => (
                  <div
                    key={step.id}
                    className={`flex flex-col items-center ${
                      index <= currentStep
                        ? 'text-blue-600 dark:text-blue-400'
                        : 'text-gray-400 dark:text-gray-600'
                    }`}
                  >
                    <div
                      className={`w-6 h-6 rounded-full flex items-center justify-center text-xs mb-1 ${
                        index <= currentStep
                          ? 'bg-blue-600 dark:bg-blue-500 text-white'
                          : 'bg-gray-200 dark:bg-gray-700 text-gray-500'
                      }`}
                    >
                      {index + 1}
                    </div>
                    <span className="text-xs whitespace-nowrap">{step.title}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Step content */}
          <div className="p-8">{renderStep()}</div>
        </div>
      </div>
    </div>
  );
};