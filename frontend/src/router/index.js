import Vue from "vue"
import Router from 'vue-router'
import Transactions from "../views/Transactions.vue";
import Settings from "../views/Settings.vue";

Vue.use(Router)

const routes = {
    transactions: {path: '/', component: Transactions},
    settings: {path: '/settings', component: Settings}
};

export {routes};
export default new Router({
    mode: 'history',
    routes: [
        routes.transactions,
        routes.settings
    ]
});