/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables

import { setupLayouts } from 'virtual:generated-layouts'
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Pages
import Index from '@/pages/index.vue'
import Login from '@/pages/login.vue'
import Signup from '@/pages/signup.vue'
import Profile from '@/pages/profile.vue'

const PUBLIC_ROUTES = new Set(['/login', '/signup', '/auth/callback'])

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'index', component: Index },
  { path: '/login', name: 'login', component: Login },
  { path: '/signup', name: 'signup', component: Signup },
  { path: '/profile', name: 'profile', component: Profile },
  { path: '/auth/callback', component: () => import('../views/OIDCCallback.vue') },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: setupLayouts(routes),
})

router.beforeEach(to => {
  const authStore = useAuthStore()
  if (!PUBLIC_ROUTES.has(to.path) && !authStore.isAuthenticated) {
    return { name: 'login' }
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

export default router
