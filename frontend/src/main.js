import Vue from 'vue';
import App from './App.vue';
import router from './router'
import Buefy, {ConfigProgrammatic} from 'buefy';
import 'buefy/dist/buefy.css'
import '@fortawesome/fontawesome-free/css/all.css'

ConfigProgrammatic.setOptions({
    defaultIconPack: 'fas'
})
Vue.use(Buefy)

new Vue({
    router: router,
    el: '#app',
    render: h => h(App),
});