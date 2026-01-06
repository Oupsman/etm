import { defineStore } from 'pinia';
import FingerprintJS from '@fingerprintjs/fingerprintjs';
const deviceStorageKey = 'etm-deviceID';
const trustedStorageKey = 'etm-deviceIsTrusted';
export const useDeviceStore = defineStore('device', {
  state: () => ({
    deviceID: localStorage.getItem(deviceStorageKey) || '',
    isTrusted: false,
    lastUsedAt: null as string | null,
  }),

  getters: {
    hasDeviceID: state => !!state.deviceID,
  },

  actions: {
    async generateDeviceID () {
      if (this.deviceID) return this.deviceID; // Déjà généré

      const fp = await FingerprintJS.load();
      const { visitorId } = await fp.get();
      this.deviceID = visitorId;
      localStorage.setItem(deviceStorageKey, visitorId);
      return visitorId;
    },

    async initDevice () {
      if (!this.deviceID) {
        await this.generateDeviceID();
      }
    },

    setTrusted (status: boolean) {
      this.isTrusted = status;
      localStorage.setItem(trustedStorageKey, String(status));
    },

    loadFromStorage () {
      const trusted = localStorage.getItem(trustedStorageKey);
      if (trusted) this.isTrusted = trusted === 'true';
    },

    // Efface le deviceID (ex: déconnexion)
    clearDevice () {
      this.deviceID = '';
      this.isTrusted = false;
      localStorage.removeItem('deviceStorageKey');
      localStorage.removeItem('deviceIsTrusted');
    },
  },
});
