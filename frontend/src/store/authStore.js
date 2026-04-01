import { create } from 'zustand';
import api from '../lib/api';

const useAuthStore = create((set, get) => ({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  isAuthenticated: !!localStorage.getItem('access_token'),
  loading: false,

  login: async (email, password) => {
    set({ loading: true });
    try {
      const res = await api.post('/auth/login', { email, password });
      const { access_token, refresh_token, user } = res.data;
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      localStorage.setItem('user', JSON.stringify(user));
      set({ user, isAuthenticated: true, loading: false });
      return { success: true };
    } catch (err) {
      set({ loading: false });
      return { success: false, error: err.response?.data?.error || 'Login failed' };
    }
  },

  register: async (name, email, password, rollNo) => {
    set({ loading: true });
    try {
      const res = await api.post('/auth/register', { name, email, password, roll_no: rollNo });
      const { access_token, refresh_token, user } = res.data;
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      localStorage.setItem('user', JSON.stringify(user));
      set({ user, isAuthenticated: true, loading: false });
      return { success: true };
    } catch (err) {
      set({ loading: false });
      return { success: false, error: err.response?.data?.error || 'Registration failed' };
    }
  },

  logout: async () => {
    try { await api.post('/auth/logout'); } catch {}
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    set({ user: null, isAuthenticated: false });
  },

  refreshUser: async () => {
    try {
      const res = await api.get('/auth/me');
      const user = res.data;
      localStorage.setItem('user', JSON.stringify(user));
      set({ user });
    } catch {}
  },

  isAdmin: () => {
    const user = get().user;
    return user && ['captain', 'vice_captain', 'manager', 'strategist'].includes(user.role);
  },

  isCaptainVC: () => {
    const user = get().user;
    return user && ['captain', 'vice_captain'].includes(user.role);
  },

  isPending: () => {
    const user = get().user;
    return user && user.join_status === 'pending';
  },

  hasTeam: () => {
    const user = get().user;
    return user && user.team_id != null;
  },
}));

export default useAuthStore;
