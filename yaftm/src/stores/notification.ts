import { ref } from 'vue'
import { defineStore } from 'pinia'
import { axiosInstance } from '@/plugins/axios'

function urlBase64ToUint8Array(base64: string): Uint8Array {
  const padding = '='.repeat((4 - (base64.length % 4)) % 4)
  const b64 = (base64 + padding).replace(/-/g, '+').replace(/_/g, '/')
  const raw = atob(b64)
  return Uint8Array.from([...raw].map(c => c.charCodeAt(0)))
}

export const useNotificationStore = defineStore('notification', () => {
  const isSupported = ref('serviceWorker' in navigator && 'PushManager' in window)
  const isSubscribed = ref(false)
  const permissionDenied = ref(false)

  async function check(): Promise<void> {
    if (!isSupported.value) return
    permissionDenied.value = Notification.permission === 'denied'
    const reg = await navigator.serviceWorker.getRegistration('/sw.js')
    if (!reg) { isSubscribed.value = false; return }
    isSubscribed.value = !!(await reg.pushManager.getSubscription())
  }

  async function subscribe(): Promise<void> {
    const { data } = await axiosInstance.get('/api/v1/getvapidkey')
    const reg = await navigator.serviceWorker.register('/sw.js', { scope: '/' })
    await navigator.serviceWorker.ready
    const sub = await reg.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(data.public_key),
    })
    await axiosInstance.post('/api/v1/user/updatesubscription', {
      subscription: JSON.stringify(sub),
    })
    isSubscribed.value = true
  }

  async function unsubscribe(): Promise<void> {
    const reg = await navigator.serviceWorker.getRegistration('/sw.js')
    if (reg) {
      const sub = await reg.pushManager.getSubscription()
      await sub?.unsubscribe()
    }
    await axiosInstance.post('/api/v1/user/updatesubscription', { subscription: '' })
    isSubscribed.value = false
  }

  async function toggle(): Promise<void> {
    if (isSubscribed.value) {
      await unsubscribe()
    } else {
      await subscribe()
    }
  }

  async function sendTest(): Promise<void> {
    await axiosInstance.post('/api/v1/user/testnotification')
  }

  return { isSupported, isSubscribed, permissionDenied, check, toggle, sendTest }
})
