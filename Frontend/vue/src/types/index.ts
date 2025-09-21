export interface LoginFormData {
  email: string;
  password: string;
}

export interface Room {
  id: number;
  name: string;
  capacity: number;
  description: string;
  features: string[];
  imageUrl?: string;
  is_available: boolean;
}

export interface User {
  id: string;
  name: string;
  email: string;
}

export interface BookingForm {
  roomId: number | null;
  date: string;
  startTime: string;
  endTime: string;
  purpose: string;
}

export interface CapacityFilter {
  value: string;
  label: string;
}

export interface BookingResponse {
  success: boolean;
  message: string;
  bookingId?: string;
}

export interface RegisterFormData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
  agreeToTerms: boolean;
}

export interface RoomForm {
  name: string;
  description: string;
  capacity: number;
  features: string[];
  is_available: boolean;
}

export interface BookingDetails {
  id: number;
  room_name: string;
  user_name: string;
  start_time: string; // ISO String date
  end_time: string; // ISO String date
  purpose: string;
  status: string;
}
