export interface Ticket {
  id: number;
  event_name: string;
  stadium: string;
  price: number;
  seat_id: string;
  status: "AVAILABLE" | "SOLD" | string;
}

export interface SearchResponse {
  data: Ticket[];
  next_cursor: string;
}

export interface BookingPayload {
  event_name: string;
  seat: string;
}

export interface BookingResponse {
  message: string;
}

export interface AuthPayload {
  email: string;
  password: string;
}

export interface AuthResponse {
  token?: string;
  message?: string;
}

export type ToastType = "success" | "error" | "info";

export interface ToastMessage {
  id: string;
  message: string;
  type: ToastType;
}