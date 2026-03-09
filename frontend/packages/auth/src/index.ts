import React, { useEffect } from 'react';
import { create } from 'zustand';

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  tenantId: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  setUser: (user: User | null) => void;
  setLoading: (loading: boolean) => void;
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: true,

  login: async (email: string, password: string) => {
    set({ isLoading: true });
    try {
      // TODO: integrate with api-client
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await response.json();
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      set({ user: data.user, isAuthenticated: true, isLoading: false });
    } catch {
      set({ isLoading: false });
      throw new Error('Login failed');
    }
  },

  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    set({ user: null, isAuthenticated: false });
  },

  setUser: (user) => set({ user, isAuthenticated: !!user }),
  setLoading: (isLoading) => set({ isLoading }),
}));

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const { setUser, setLoading } = useAuth();

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token) {
      // TODO: validate token and fetch user profile
      setLoading(false);
    } else {
      setLoading(false);
    }
  }, [setUser, setLoading]);

  return children;
}

interface ProtectedRouteProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function ProtectedRoute({ children, fallback }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return null;
  }

  if (!isAuthenticated) {
    return fallback ?? null;
  }

  return children;
}
