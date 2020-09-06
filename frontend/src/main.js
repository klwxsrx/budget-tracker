import Vue from 'vue';
import Buefy, {ConfigProgrammatic} from 'buefy';
import 'buefy/dist/buefy.css'
import '@fortawesome/fontawesome-free/css/all.css'
import App from './components/App.vue';

ConfigProgrammatic.setOptions({
    defaultIconPack: 'fas'
})
Vue.use(Buefy)
new Vue({
    el: '#app',
    render: h => h(App),
});