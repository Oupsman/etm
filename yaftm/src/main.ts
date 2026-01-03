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

const app = createApp(App)

const emitter = mitt()

registerPlugins(app)
app.config.globalProperties.emitter = emitter

app.mount('#app')
