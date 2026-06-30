/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'
import mitt from 'mitt'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)

const emitter = mitt()

registerPlugins(app)
app.config.globalProperties.emitter = emitter

app.mount('#app')

const authStore = useAuthStore()
authStore.initRefreshSchedule()

document.addEventListener('visibilitychange', () => {
  if (document.visibilityState === 'visible') authStore.initRefreshSchedule()
})
