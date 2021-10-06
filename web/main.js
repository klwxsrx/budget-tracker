import Vue from 'vue'
import {getVueI18n} from './i18n/I18nService'
import VueMq from 'vue-mq'
import App from './App.vue'
import store from './store'
import router from './router'
import Buefy from 'buefy' // TODO: get specific components
import {library} from '@fortawesome/fontawesome-svg-core'
import {faMoneyBillAlt} from '@fortawesome/free-solid-svg-icons'
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome'
import './scss/variables.scss'
import './scss/buefy.scss'

library.add(faMoneyBillAlt)
Vue.component('FontAwesomeIcon', FontAwesomeIcon)
Vue.use(Buefy, {
  defaultIconComponent: 'FontAwesomeIcon',
  defaultIconPack: 'fas',
})

Vue.use(VueMq, {
  breakpoints: {
    tablet: 769, // coupling with bulma/sass/utilities/initial-variables.sass: $tablet
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