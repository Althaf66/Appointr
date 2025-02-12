export interface Mentor {
  id: string;
  name: string;
  role: string;
  company: string;
  avatar: string;
  expertise: string[];
  availability: string;
}

export interface Session {
  id: string;
  title: string;
  mentor: Mentor;
  date: string;
  duration: string;
  spots: number;
}