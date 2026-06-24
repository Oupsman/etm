import type { Ref } from 'vue'
import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import router from '@/router'
import axios from 'axios'
import { useSnackbarStore } from '@/stores/snackbar'
import { useAuthStore } from '@/stores/auth'
import { axiosInstance } from '@/plugins/axios'
import type { User, UserSession } from '@/types/user'
function parseJwt (token: string) {
  const base64Url = token.split('.')[1]
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
  const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
  }).join(''))
  return JSON.parse(jsonPayload)
}

export const useUserStore = defineStore('user', () => {
  const snackbar = useSnackbarStore()
  const session: Ref<UserSession | null> = ref(null)

  const login = async (username: string, password: string): Promise<void> => {
    const authStore = useAuthStore()

    try {
      const response = await axios.post('/api/v1/user/login', { username, password })
      if (response.data.token) {
        authStore.setTokens(response.data.token, '')
        setUserSession({ ...parseJwt(response.data.token), token: response.data.token })
        router.push('/')
      }
    } catch (error) {
      snackbar.showSnackbar({ message: 'Login failed! ' + (error instanceof Error ? error.message : String(error)), color: 'error' })
      throw new Error('Login failed')
    }
  }

  const logout = async (): Promise<void> => {
    const authStore = useAuthStore()
    try {
      await axiosInstance.get('/api/v1/user/logout')
    } catch {
      // ignore logout errors
    }
    authStore.logout()
    session.value = null
    router.push('/login')
  }

  const signup = async (username: string, password: string, email: string): Promise<void> => {
    try {
      await axios.post('/api/v1/user/register', { username, password, email })
      snackbar.showSnackbar({ message: 'Signup successful', color: 'success' })
      router.push('/login')
    } catch (error) {
      snackbar.showSnackbar({ message: 'Signup failed! ' + (error instanceof Error ? error.message : String(error)), color: 'error' })
      throw new Error('Signup failed')
    }
  }

  const setUserSession = (data: UserSession): void => {
    session.value = data
  }

  const userIsLoggedIn = computed(() => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (!session.value && token) {
      session.value = parseJwt(token)
    }
    if (session.value?.exp) {
      const expiresAt = new Date(0).setUTCSeconds(session.value.exp)
      return Date.now() < expiresAt
    }
    return false
  })

  const checkToken = computed((): boolean => {
    const authStore = useAuthStore()
    const token = authStore.token
    if (token) {
      const decoded = parseJwt(token)
      if (decoded.exp < Date.now() / 1000) {
        authStore.logout()
        session.value = null
        return false
      }
      session.value = decoded
    }
    return true
  })

  const user = ref({})

  const getUser = async (): Promise<User> => {
    const response = await axiosInstance.get('/api/v1/user')
    user.value = response.data
    return response.data
  }

  const updateUser = async (userData: object): Promise<boolean> => {
    await axiosInstance.post('/api/v1/user', userData)
    return true
  }

  const saveActiveCategory = async (categoryID: number): Promise<void> => {
    await axiosInstance.patch('/api/v1/user/preferences', { active_category_id: categoryID })
  }

  const trustUserDevice = async (userData: object): Promise<boolean> => {
    await axiosInstance.post('/api/v1/user/devices', userData)
    return true
  }

  return { session, userIsLoggedIn, login, logout, setUserSession, checkToken, signup, user, getUser, updateUser, trustUserDevice, saveActiveCategory }
})
