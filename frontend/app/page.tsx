"use client"; // This component needs to run in the browser

import { useEffect, useState } from "react";
import { Ticket } from "../src/types";
import TicketCard from "../src/components/TicketCard";

export default function Home() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);

  // Function to fetch data
  const fetchTickets = async () => {
    try {
      // Connect to Go API (Port 8081)
      const res = await fetch("http://localhost:8081/search?q=Concert&limit=20");
      const data = await res.json();
      
      // The API returns { data: [...], next_cursor: "..." }
      console.log(data)
      if (data.data) {
        setTickets(data.data);
      }
    } catch (error) {
      console.error("Failed to fetch tickets:", error);
    } finally {
      setLoading(false);
    }
  };

  // Run once when page loads
  useEffect(() => {
    fetchTickets();
  }, []);

  return (
    <main className="min-h-screen bg-black text-white p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-600 mb-8 text-center">
          🎟️ TicketRush Live
        </h1>

        {loading ? (
          <div className="text-center animate-pulse">Scanning the blockchain...</div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {tickets.map((ticket) => (
              <TicketCard key={ticket.id} ticket={ticket} />
            ))}
          </div>
        )}
      </div>
    </main>
  );
}