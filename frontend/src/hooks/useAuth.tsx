"use client";

import { useCallback, useState } from "react";
import { authApi } from "../services/api";
import { decodeJwt } from "../lib/utils";
import type { AuthPayload } from "../types";

interface JwtPayload {
  sub?: string;
  email?: string;
  exp?: number;
}

interface AuthState {
  token: string | null;
  email: string | null;
}

function loadFromSession(): AuthState {
  if (typeof window === "undefined") return { token: null, email: null };
  const token = sessionStorage.getItem("token");
  if (!token) return { token: null, email: null };
  const payload = decodeJwt<JwtPayload>(token);
  if (!payload || (payload.exp && payload.exp * 1000 < Date.now())) {
    sessionStorage.removeItem("token");
    return { token: null, email: null };
  }
  return { token, email: payload.email ?? payload.sub ?? null };
}

export function useAuth() {
  const [auth, setAuth] = useState<AuthState>(loadFromSession);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = useCallback(async (payload: AuthPayload) => {
    setLoading(true);
    setError(null);
    try {
      const res = await authApi.login(payload);
      if (!res.token) throw new Error("No token received");
      const decoded = decodeJwt<JwtPayload>(res.token);
      sessionStorage.setItem("token", res.token);
      setAuth({ token: res.token, email: decoded?.email ?? decoded?.sub ?? null });
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const register = useCallback(async (payload: AuthPayload) => {
    setLoading(true);
    setError(null);
    try {
      await authApi.register(payload);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(() => {
    sessionStorage.removeItem("token");
    setAuth({ token: null, email: null });
  }, []);

  return {
    token: auth.token,
    email: auth.email,
    isLoggedIn: !!auth.token,
    loading,
    error,
    login,
    register,
    logout,
    clearError: () => setError(null),
  };
}