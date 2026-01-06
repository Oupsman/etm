import { defineStore } from 'pinia';
import axios from 'axios';
import { useDeviceStore } from './device';

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
    },

    logout () {
      this.token = '';
      this.refreshToken = '';
      this.userID = '';
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('userID');
      const deviceStore = useDeviceStore();
      deviceStore.clearDevice();
    },

    async refreshAccessToken () {
      try {
        const deviceStore = useDeviceStore();
        const { data } = await axios.post('/refresh-token', {
          refreshToken: this.refreshToken,
        }, {
          headers: { 'X-Device-ID': deviceStore.deviceID },
        });
        this.setTokens(data.token, data.refreshToken);
        return data.token;
      } catch (error) {
        this.logout();
        throw error;
      }
    },
  },
});
