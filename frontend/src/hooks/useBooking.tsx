"use client";

import { useCallback, useState } from "react";
import { bookingApi } from "../services/api";
import type { BookingPayload } from "../types";

export function useBooking(token: string | null) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const book = useCallback(
    async (payload: BookingPayload) => {
      if (!token) {
        setError("You must be logged in to book a ticket.");
        return false;
      }
      setLoading(true);
      setError(null);
      setSuccess(null);
      try {
        const res = await bookingApi.book(payload, token);
        setSuccess(res.message ?? "Booking confirmed!");
        return true;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Booking failed");
        return false;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  const reset = useCallback(() => {
    setError(null);
    setSuccess(null);
  }, []);

  return { loading, error, success, book, reset };
}