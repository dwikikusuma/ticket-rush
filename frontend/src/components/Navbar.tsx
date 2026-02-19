// src/components/Navbar.tsx
"use client";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="bg-white border-b border-gray-100 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <Link href="/" className="text-2xl font-black text-blue-600 tracking-tighter">
            TicketRush<span className="text-black">.</span>
          </Link>
          
          <div className="flex gap-6 items-center">
            <Link href="/" className="text-gray-600 hover:text-black font-medium text-sm">Browse</Link>
            {user ? (
              <>
                <Link href="/my-tickets" className="text-gray-600 hover:text-black font-medium text-sm">My Tickets</Link>
                <button 
                  onClick={logout}
                  className="bg-gray-100 text-gray-900 px-4 py-2 rounded-lg text-sm font-semibold hover:bg-gray-200 transition"
                >
                  Logout
                </button>
              </>
            ) : (
              <Link 
                href="/login" 
                className="bg-blue-600 text-white px-5 py-2 rounded-lg text-sm font-semibold hover:bg-blue-700 transition shadow-sm"
              >
                Sign In
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}