import { computed, ref } from 'vue'
import type { Ref } from 'vue'
import { defineStore } from 'pinia'
import router from '@/router'
import axios from 'axios'

import type { UserSession } from '@/types/user'

function parseJwt (token: string) {
  const base64Url = token.split('.')[1]
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
  const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
  }).join(''))

  return JSON.parse(jsonPayload)
}

export const useUserStore = defineStore('user', () => {
  const session: Ref<UserSession | null> = ref(null)

  const login = async (username: string, password:string, callback?: Function): Promise<void> => {
    axios.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/user/login', {
      username,
      password,
    }).then(response => {
      if (response.data.token) {
        localStorage.setItem('etm-token', response.data.token)
        const data = parseJwt(response.data.token)
        setUserSession({
          ...data,
          token: response.data.token,
        })
        router.push('/')
      }
    }).catch(error => {
      console.error('Login error:', error)
      throw new Error('Login failed')
    })
  }

  const logout = async (): Promise<void> => {
    const token = localStorage.getItem('etm-token')
    axios.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/user/logout', {
      token,
    }).then(response => {
      if (response.data.token) {
        localStorage.removeItem('etm-token')
        session.value = null
        router.push('/')
      }
    }).catch(error => {
      console.error('Login error:', error)
      throw new Error('Login failed')
    })
  }
  const signup = async (username: string, password: string, email: string): Promise<void> => {
    const request = axios.create(
      {
        baseURL: import.meta.env.VITE_BACKEND_URL,
        timeout: 1000,
        withCredentials: false,
      }
    )
    request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/user/register', {
      username,
      password,
      email,
    }).then(response => {
      console.log('Signup response:', response)
    }).catch(error => {
      console.error('Signup error:', error)
      throw new Error('Signup failed')
    })
  }

  const setUserSession = (data: UserSession | any): void => {
    session.value = data
  }

  const userIsLoggedIn = computed(() => {
    if (session.value?.exp) {
      const expiresAt = new Date(0).setUTCSeconds(session.value.exp)
      const now = new Date().getSeconds()
      return now < expiresAt
    }
    return false
  })
  // checks if the JWT token is still valid
  const checkToken = computed((): boolean => {
    const token = localStorage.getItem('etm-token')
    if (token) {
      const decoded = parseJwt(token)
      if (decoded.exp < Date.now() / 1000) {
        // Token is expired
        localStorage.removeItem('etm-token')
        session.value = null
        return false
      } else {
        session.value = decoded
      }
    }
    return true
  })

  const user = ref({})

  const getUser = async () : Promise<Object> => {
    console.log('Get user - function')
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }

    try {
      const request = axios.create({
        baseURL: import.meta.env.VITE_BACKEND_URL,
        timeout: 1000,
        headers: { Authorization: `Bearer ${token}` },
      })

      const response = await request.get(`/api/v1/user`)
      console.log('Response : ', response.data)
      user.value = response.data
      console.log('user:', response.data)
      return response.data.activity
    } catch (error) {
      console.error('Get user error:', error)
      throw new Error('get user failed')
    }
  }

  const updateUser = async (user: Object) : Promise<Object> => {
    console.log('user to save', user)
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 30000,
      headers: { Authorization: `Bearer ${token}` },
    })
    request.post(import.meta.env.VITE_BACKEND_URL + '/api/v1/user', {
      ...user,
    }).then(response => {
      console.log('User updated:', response.data)
      return true
    }).catch(error => {
      console.error('Update user error:', error)
      throw new Error('update user')
    })
    return true
  }

  const refreshDashboard = async (): Promise<void> => {
    console.log('Get dashboard - function')
    const token = localStorage.getItem('etm-token')
    if (!token) {
      throw new Error('No token')
    }
    const request = axios.create({
      baseURL: import.meta.env.VITE_BACKEND_URL,
      timeout: 30000,
      headers: { Authorization: `Bearer ${token}` },
    })

  }
  return { session, userIsLoggedIn, login, logout, setUserSession, checkToken, signup, refreshDashboard, user, getUser, updateUser }
})
