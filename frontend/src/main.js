import Vue from 'vue'
import {getVueI18n} from './services/I18nService'
import App from './App.vue'
import store from './store'
import router from './router'
import Buefy, {ConfigProgrammatic} from 'buefy'
import './scss/variables.scss'
import './scss/buefy.scss'

ConfigProgrammatic.setOptions({
  defaultIconPack: 'fas',
})
Vue.use(Buefy)

new Vue({
  i18n: getVueI18n(),
  store: store,
  router: router,
  el: '#app',
  render: h => h(App),
})