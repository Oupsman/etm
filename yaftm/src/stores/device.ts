import { defineStore } from 'pinia'
import { axiosInstance } from '@/plugins/axios'
import type { Device } from '@/types/user'

export const useDeviceStore = defineStore('device', {
  state: () => ({
    devices: [] as Device[],
    loading: false,
  }),

  actions: {
    async fetchDevices () {
      this.loading = true
      try {
        const { data } = await axiosInstance.get('/api/v1/devices')
        this.devices = data.devices ?? []
      } finally {
        this.loading = false
      }
    },

    async registerDevice (name: string): Promise<string> {
      const { data } = await axiosInstance.post('/api/v1/devices', { device_name: name })
      this.devices.push(data.device)
      return data.api_key as string
    },

    async deleteDevice (id: number) {
      await axiosInstance.delete(`/api/v1/devices/${id}`)
      this.devices = this.devices.filter(d => d.ID !== id)
    },
  },
})
