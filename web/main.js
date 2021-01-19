import Vue from 'vue'
import {getVueI18n} from './services/I18nService'
import VueMq from 'vue-mq'
import App from './App.vue'
import store from './store'
import router from './router'
import Buefy, {ConfigProgrammatic} from 'buefy'
import './scss/variables.scss'
import './scss/buefy.scss'
import '@fortawesome/fontawesome-free/css/all.css'

ConfigProgrammatic.setOptions({
  defaultIconPack: 'fas',
})
Vue.use(Buefy)
Vue.use(VueMq, {
  breakpoints: {
    tablet: 769, // coupling with scss/variables.scss: $tablet
    desktop: Infinity,
  },
})

new Vue({
  i18n: getVueI18n(),
  store: store,
  router: router,
  el: '#app',
  render: h => h(App),
})