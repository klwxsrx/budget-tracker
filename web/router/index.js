import Vue from 'vue'
import Router from 'vue-router'
import TransactionsPage from '../view/TransactionsPage.vue'
import SettingsPage from '../view/SettingsPage.vue'

Vue.use(Router)

const routes = {
  transactions: {path: '/', component: TransactionsPage},
  settings: {path: '/settings', component: SettingsPage},
}

export {routes}
export default new Router({
  mode: 'history',
  routes: [
    routes.transactions,
    routes.settings,
  ],
})