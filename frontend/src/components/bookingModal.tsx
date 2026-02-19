"use client";

import { useState } from "react";
import type { Ticket } from "@/types";
import { formatPrice } from "@/lib/utils";

interface Props {
  ticket: Ticket;
  onConfirm: (eventName: string, seat: string) => Promise<boolean>;
  onClose: () => void;
  loading: boolean;
  error: string | null;
  success: string | null;
}

export default function BookingModal({ ticket, onConfirm, onClose, loading, error, success }: Props) {
  const [seat, setSeat] = useState(ticket.seat_id);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onConfirm(ticket.event_name, seat);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm">
      <div className="bg-gray-900 border border-gray-700 rounded-xl shadow-2xl w-full max-w-md mx-4">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-800">
          <h2 className="text-lg font-bold text-white">Confirm Booking</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white transition text-xl leading-none">×</button>
        </div>

        <div className="p-6">
          {success ? (
            <div className="text-center py-6">
              <p className="text-4xl mb-3">🎟️</p>
              <p className="text-emerald-400 font-semibold text-lg">{success}</p>
              <button
                onClick={onClose}
                className="mt-6 bg-gray-800 hover:bg-gray-700 text-white px-6 py-2 rounded-lg transition"
              >
                Close
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Ticket info */}
              <div className="bg-gray-800 rounded-lg p-4 space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-400">Event</span>
                  <span className="text-white font-medium">{ticket.event_name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Venue</span>
                  <span className="text-white">{ticket.stadium}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Price</span>
                  <span className="text-green-400 font-bold">{formatPrice(ticket.price)}</span>
                </div>
              </div>

              <div>
                <label className="block text-gray-400 text-sm mb-1">Seat</label>
                <input
                  type="text"
                  value={seat}
                  onChange={(e) => setSeat(e.target.value)}
                  required
                  className="w-full bg-gray-800 border border-gray-700 text-white rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {error && <p className="text-red-400 text-sm">⚠️ {error}</p>}

              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={onClose}
                  className="flex-1 bg-gray-800 hover:bg-gray-700 text-white py-2.5 rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg transition"
                >
                  {loading ? "Booking…" : "Confirm"}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}