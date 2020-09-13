import Vue from 'vue'
import {getVueI18n} from './services/I18nService'
import App from './App.vue'
import router from './router'
import Buefy, {ConfigProgrammatic} from 'buefy'
import 'buefy/dist/buefy.css'
import '@fortawesome/fontawesome-free/css/all.css'

ConfigProgrammatic.setOptions({
    defaultIconPack: 'fas',
})

Vue.use(Buefy)

new Vue({
    i18n: getVueI18n(),
    router: router,
    el: '#app',
    render: h => h(App),
})