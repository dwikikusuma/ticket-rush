"use client";

import { useEffect, useState } from "react";
import TicketCard from "@/components/ticketCard";
import SearchBar from "@/components/searchBar";
import SkeletonCard from "@/components/skeletonCard";
import StatsBar from "@/components/statsBar";
import Toast from "@/components/toast";
import AuthModal from "@/components/authModal";
import BookingModal from "@/components/bookingModal";
import { useAuth } from "@/hooks/useAuth";
import { useTickets } from "@/hooks/useTickets";
import { useBooking } from "@/hooks/useBooking";
import { useToast } from "@/hooks/useToast";
import type { Ticket } from "@/types";

const SKELETON_COUNT = 6;

export default function Home() {
  const auth = useAuth();
  const tickets = useTickets();
  const booking = useBooking(auth.token);
  const { toasts, addToast, dismiss } = useToast();

  const [showAuthModal, setShowAuthModal] = useState(false);
  const [bookingTarget, setBookingTarget] = useState<Ticket | null>(null);

  // Initial load
  useEffect(() => {
    tickets.search("Concert");
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Toast on booking success
  useEffect(() => {
    if (booking.success) addToast(booking.success, "success");
  }, [booking.success, addToast]);

  const handleBook = (ticket: Ticket) => {
    if (!auth.isLoggedIn) {
      setShowAuthModal(true);
      return;
    }
    booking.reset();
    setBookingTarget(ticket);
  };

  const handleConfirmBooking = async (eventName: string, seat: string) => {
    const ok = await booking.book({ event_name: eventName, seat });
    if (!ok && booking.error) addToast(booking.error, "error");
    return ok;
  };

  const handleCloseBooking = () => {
    setBookingTarget(null);
    booking.reset();
  };

  const handleLogin = async (email: string, password: string) => {
    const ok = await auth.login({ email, password });
    if (ok) {
      setShowAuthModal(false);
      addToast("Logged in successfully!", "success");
    }
    return ok;
  };

  const handleRegister = async (email: string, password: string) => {
    return auth.register({ email, password });
  };

  return (
    <main className="min-h-screen bg-black text-white p-6 md:p-10">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-500">
            🎟️ TicketRush
          </h1>
          <div className="flex items-center gap-3">
            {auth.isLoggedIn ? (
              <>
                <span className="text-gray-400 text-sm hidden sm:block">{auth.email}</span>
                <button
                  onClick={auth.logout}
                  className="text-sm text-gray-400 hover:text-white transition"
                >
                  Log out
                </button>
              </>
            ) : (
              <button
                onClick={() => setShowAuthModal(true)}
                className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-4 py-2 rounded-lg transition"
              >
                Log in
              </button>
            )}
          </div>
        </div>

        {/* Search */}
        <div className="mb-6">
          <SearchBar
            value={tickets.query}
            onChange={tickets.onQueryChange}
            placeholder="Search events, venues…"
          />
        </div>

        {/* Stats */}
        <div className="mb-6">
          <StatsBar tickets={tickets.tickets} />
        </div>

        {/* Error */}
        {tickets.error && (
          <p className="text-red-400 text-sm mb-4">⚠️ {tickets.error}</p>
        )}

        {/* Grid */}
        {tickets.loading ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: SKELETON_COUNT }).map((_, i) => (
              <SkeletonCard key={i} />
            ))}
          </div>
        ) : tickets.tickets.length === 0 ? (
          <div className="text-center text-gray-500 py-20">
            <p className="text-4xl mb-3">🎭</p>
            <p>No tickets found. Try a different search.</p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {tickets.tickets.map((ticket) => (
              <TicketCard
                key={ticket.id}
                ticket={ticket}
                onBook={handleBook}
                isLoggedIn={auth.isLoggedIn}
              />
            ))}
          </div>
        )}

        {/* Load More */}
        {tickets.hasMore && !tickets.loading && (
          <div className="mt-8 text-center">
            <button
              onClick={tickets.loadMore}
              disabled={tickets.loadingMore}
              className="bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white px-6 py-2.5
                         rounded-lg transition"
            >
              {tickets.loadingMore ? "Loading…" : "Load more"}
            </button>
          </div>
        )}
      </div>

      {/* Modals */}
      {showAuthModal && (
        <AuthModal
          onLogin={handleLogin}
          onRegister={handleRegister}
          onClose={() => { setShowAuthModal(false); auth.clearError(); }}
          loading={auth.loading}
          error={auth.error}
        />
      )}

      {bookingTarget && (
        <BookingModal
          ticket={bookingTarget}
          onConfirm={handleConfirmBooking}
          onClose={handleCloseBooking}
          loading={booking.loading}
          error={booking.error}
          success={booking.success}
        />
      )}

      {/* Toasts */}
      <Toast toasts={toasts} onDismiss={dismiss} />
    </main>
  );
}