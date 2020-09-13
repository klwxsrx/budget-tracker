import Vue from 'vue';
import VueI18n from 'vue-i18n'
import App from './App.vue';
import router from './router'
import Buefy, {ConfigProgrammatic} from 'buefy';
import 'buefy/dist/buefy.css'
import '@fortawesome/fontawesome-free/css/all.css'

ConfigProgrammatic.setOptions({
    defaultIconPack: 'fas'
})

Vue.use(VueI18n)
Vue.use(Buefy)

const supportedLocales = [
    'en',
    'ru'
];
const userLocale = (navigator.userLanguage || navigator.language).slice(0, 2);
const i18n = new VueI18n({
    locale: supportedLocales.includes(userLocale) ? userLocale : supportedLocales[0],
})

new Vue({
    i18n: i18n,
    router: router,
    el: '#app',
    render: h => h(App),
});