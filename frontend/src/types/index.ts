export interface Mentor {
  id: string;
  userid: string;
  name: string;
  avatar?: string;
  country?: string;
  hourlyRate?: number;
  availabilityDetails?: string;
  language?: string[];
  education?: { degree: string; institution: string; year: string }[];
  gigs?: { title: string; description: string; expertise: string; discipline: string[] }[];
  experience?: { title: string; company: string; description: string; year_from: string; year_to: string }[];
  workingat?: { title: string; company: string; totalyear: number };
  bookingslots?: { days: string[]; start_time: string; end_time: string; start_period: string; end_period: string }[];
}

export interface Session {
  id: string;
  title: string;
  mentor: Mentor;
  date: string;
  duration: string;
  spots: number;
}