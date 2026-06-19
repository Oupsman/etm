import { defineStore } from 'pinia'
import { ref } from 'vue'
import i18n from '@/plugins/i18n'

export const useLocaleStore = defineStore('locale', () => {
  const locale = ref(localStorage.getItem('locale') ?? 'en')

  const setLocale = (newLocale: string) => {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
    i18n.global.locale.value = newLocale as 'en' | 'fr'
  }

  return { locale, setLocale }
})
