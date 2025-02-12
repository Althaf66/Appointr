export interface MentorProfile {
  id: string;
  name: string;
  role: string;
  company: string;
  avatar: string;
  location: string;
  bio: string;
  expertise: string[];
  languages: string[];
  education: Education[];
  experience: Experience[];
  reviews: Review[];
  sessions: Session[];
}

export interface Education {
  school: string;
  degree: string;
  field: string;
  year: string;
}

export interface Experience {
  company: string;
  role: string;
  duration: string;
  description: string;
}

export interface Review {
  id: string;
  author: string;
  avatar: string;
  rating: number;
  content: string;
  date: string;
}

export interface Session {
  id: string;
  title: string;
  duration: string;
  spots: number;
  price: string;
  date: string;
}