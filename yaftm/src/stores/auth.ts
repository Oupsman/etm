import { defineStore } from 'pinia';
import axios from 'axios';

let refreshTimer: ReturnType<typeof setTimeout> | null = null;

function getTokenExp (token: string): number | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return typeof payload.exp === 'number' ? payload.exp : null;
  } catch {
    return null;
  }
}

function scheduleRefresh (token: string, refreshFn: () => Promise<void>): void {
  if (refreshTimer) clearTimeout(refreshTimer);
  const exp = getTokenExp(token);
  if (!exp) return;
  const msUntilRefresh = (exp * 1000) - Date.now() - 5 * 60 * 1000;
  if (msUntilRefresh <= 0) {
    refreshFn();
    return;
  }
  refreshTimer = setTimeout(refreshFn, msUntilRefresh);
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    refreshToken: localStorage.getItem('refreshToken') || '',
    userID: localStorage.getItem('userID') || '',
  }),

  getters: {
    isAuthenticated: state => !!state.token,
  },

  actions: {
    setTokens (token: string, refreshToken: string) {
      this.token = token;
      this.refreshToken = refreshToken;
      localStorage.setItem('token', token);
      localStorage.setItem('refreshToken', refreshToken);
      scheduleRefresh(token, () => this.refreshAccessToken());
    },

    logout () {
      if (refreshTimer) { clearTimeout(refreshTimer); refreshTimer = null; }
      this.token = '';
      this.refreshToken = '';
      this.userID = '';
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('userID');
    },

    async refreshAccessToken () {
      try {
        const { data } = await axios.get('/api/v1/user/refreshtoken', {
          headers: { Authorization: `Bearer ${this.token}` },
        });
        if (!data.token) throw new Error('no token in response');
        this.setTokens(data.token, this.refreshToken);
        return data.token;
      } catch (error) {
        this.logout();
        throw error;
      }
    },

    initRefreshSchedule () {
      if (this.token) scheduleRefresh(this.token, () => this.refreshAccessToken());
    },
  },
});
