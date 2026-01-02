/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
// import { registerPlugins } from '@/plugins'
import mitt from 'mitt'
import { createPinia } from 'pinia'
import { createVuetify } from 'vuetify'

// Components
import App from './App.vue'
import router from './router'

// Composables
import { createApp } from 'vue'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const app = createApp(App)

const emitter = mitt()
// Handled by registerPlugins

const pinia = createPinia()
const vuetify = createVuetify({
  components,
  directives,
})

app.use(pinia)
app.use(router)
app.use(vuetify)

// registerPlugins(app)
app.config.globalProperties.emitter = emitter

app.mount('#app')
