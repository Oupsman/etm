// stores/snackbar.ts
import { defineStore } from 'pinia';

type SnackbarColor = 'info' | 'success' | 'error' | 'warning' | 'primary';

interface SnackbarState {
  show: boolean;
  message: string;
  color: SnackbarColor;
  timeout: number;
  actionText?: string;
  actionCallback?: () => void;
}

export const useSnackbarStore = defineStore('snackbar', {
  state: (): SnackbarState => ({
    show: false,
    message: '',
    color: 'info',
    timeout: 3000,
  }),
  actions: {
    showSnackbar (payload: {
      message: string;
      color?: SnackbarColor;
      timeout?: number;
      actionText?: string;
      actionCallback?: () => void;
    }) {
      this.message = payload.message;
      this.color = payload.color || 'info';
      this.timeout = payload.timeout || 3000;
      this.actionText = payload.actionText;
      this.actionCallback = payload.actionCallback;
      this.show = true;
    },
    closeSnackbar () {
      this.show = false;
      this.actionText = undefined;
      this.actionCallback = undefined;
    },
  },
});
