// plugins/axios.ts
import axios from 'axios';
import { useAuthStore } from '@/stores/auth';
import router from '@/router';

export const axiosInstance = axios.create();

axiosInstance.interceptors.request.use(config => {
  const authStore = useAuthStore();
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`;
  }
  return config;
});

axiosInstance.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;
    const authStore = useAuthStore();

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const newToken = await authStore.refreshAccessToken();
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return axiosInstance(originalRequest);
      } catch {
        authStore.logout();
        router.push('/login');
        return Promise.reject(error);
      }
    }
    return Promise.reject(error);
  }
);
