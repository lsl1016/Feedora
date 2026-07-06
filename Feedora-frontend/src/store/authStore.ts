import { create } from 'zustand';
import type { User } from '../types';

interface AuthState {
  token: string | null;
  currentUser: User | null;
  isLogin: boolean;
  setAuth: (token: string, user: User) => void;
  setCurrentUser: (user: User | null) => void;
  logout: () => void;
}

const rawToken = localStorage.getItem('community_v21_token');
const rawUser = localStorage.getItem('community_v21_user');

export const useAuthStore = create<AuthState>((set) => ({
  token: rawToken,
  currentUser: rawUser ? JSON.parse(rawUser) : null,
  isLogin: Boolean(rawToken),
  setAuth: (token, user) => {
    localStorage.setItem('community_v21_token', token);
    localStorage.setItem('community_v21_user', JSON.stringify(user));
    set({ token, currentUser: user, isLogin: true });
  },
  setCurrentUser: (user) => {
    if (user) localStorage.setItem('community_v21_user', JSON.stringify(user));
    set({ currentUser: user });
  },
  logout: () => {
    localStorage.removeItem('community_v21_token');
    localStorage.removeItem('community_v21_user');
    set({ token: null, currentUser: null, isLogin: false });
  },
}));
