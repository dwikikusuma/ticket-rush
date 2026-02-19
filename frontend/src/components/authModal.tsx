"use client";

import { useState } from "react";

interface Props {
  onLogin: (email: string, password: string) => Promise<boolean>;
  onRegister: (email: string, password: string) => Promise<boolean>;
  onClose: () => void;
  loading: boolean;
  error: string | null;
}

export default function AuthModal({ onLogin, onRegister, onClose, loading, error }: Props) {
  const [tab, setTab] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [regSuccess, setRegSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (tab === "login") {
      await onLogin(email, password);
    } else {
      const ok = await onRegister(email, password);
      if (ok) setRegSuccess(true);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm">
      <div className="bg-gray-900 border border-gray-700 rounded-xl shadow-2xl w-full max-w-md mx-4">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-800">
          <h2 className="text-lg font-bold text-white">
            {tab === "login" ? "Sign In" : "Create Account"}
          </h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white transition text-xl leading-none">×</button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-gray-800">
          {(["login", "register"] as const).map((t) => (
            <button
              key={t}
              onClick={() => { setTab(t); setRegSuccess(false); }}
              className={`flex-1 py-3 text-sm font-medium transition ${
                tab === t ? "text-blue-400 border-b-2 border-blue-400" : "text-gray-500 hover:text-gray-300"
              }`}
            >
              {t === "login" ? "Login" : "Register"}
            </button>
          ))}
        </div>

        <div className="p-6">
          {regSuccess ? (
            <div className="text-center py-4">
              <p className="text-emerald-400 font-medium">✅ Account created!</p>
              <p className="text-gray-400 text-sm mt-1">You can now log in.</p>
              <button onClick={() => { setTab("login"); setRegSuccess(false); }} className="mt-4 text-blue-400 text-sm hover:underline">
                Go to Login →
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-gray-400 text-sm mb-1">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="w-full bg-gray-800 border border-gray-700 text-white rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="you@example.com"
                />
              </div>
              <div>
                <label className="block text-gray-400 text-sm mb-1">Password</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="w-full bg-gray-800 border border-gray-700 text-white rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="••••••••"
                />
              </div>

              {error && <p className="text-red-400 text-sm">⚠️ {error}</p>}

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg transition"
              >
                {loading ? "Please wait…" : tab === "login" ? "Sign In" : "Create Account"}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}