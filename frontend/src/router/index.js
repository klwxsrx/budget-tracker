import Vue from "vue"
import Router from 'vue-router'
import Transactions from "../views/Transactions.vue";
import Settings from "../views/Settings.vue";

Vue.use(Router)

export default new Router({
    mode: 'history',
    routes: [
        {path: '/', component: Transactions},
        {path: '/settings', component: Settings}
    ]
});