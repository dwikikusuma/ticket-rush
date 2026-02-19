import type {
  AuthPayload,
  AuthResponse,
  BookingPayload,
  BookingResponse,
  SearchResponse,
} from "../types";

const AUTH_URL = "http://localhost:8087";
const SEARCH_URL = "http://localhost:8083";
const BOOKING_URL = "http://localhost:8086";

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText);
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export const authApi = {
  login: (payload: AuthPayload): Promise<AuthResponse> =>
    fetch(`${AUTH_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    }).then(handleResponse<AuthResponse>),

  register: (payload: AuthPayload): Promise<AuthResponse> =>
    fetch(`${AUTH_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    }).then(handleResponse<AuthResponse>),
};

export const searchApi = {
  search: (
    q: string,
    limit = 20,
    cursor?: string,
    token?: string | null
  ): Promise<SearchResponse> => {
    const params = new URLSearchParams({ q, limit: String(limit) });
    if (cursor) params.set("cursor", cursor);
    const headers: Record<string, string> = {};
    if (token) headers["Authorization"] = `Bearer ${token}`;
    return fetch(`${SEARCH_URL}/search?${params}`, { headers }).then(
      handleResponse<SearchResponse>
    );
  },
};

export const bookingApi = {
  book: (payload: BookingPayload, token: string): Promise<BookingResponse> =>
    fetch(`${BOOKING_URL}/bookings`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(payload),
    }).then(handleResponse<BookingResponse>),
};