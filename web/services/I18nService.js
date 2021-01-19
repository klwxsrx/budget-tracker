import Vue from 'vue'
import VueI18n from 'vue-i18n'
import moment from 'moment'

const supportedLocales = [
  'en',
  'ru',
]

Vue.use(VueI18n)
let vueI18n = new VueI18n({
  locale: getUserLocale(),
})

moment.locale(vueI18n.locale)

function getUserLocale() {
  const userLocale = (navigator.userLanguage || navigator.language).slice(0, 2)
  return supportedLocales.includes(userLocale) ? userLocale : supportedLocales[0]
}

export function getVueI18n() {
  return vueI18n
}

export default {
  getCurrentLocale() {
    return getVueI18n().locale
  },
}