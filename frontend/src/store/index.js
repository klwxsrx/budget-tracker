import Vue from 'vue'
import Vuex from 'vuex'
import account from './modules/account'
import settings from './modules/settings'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    account,
    settings,
  },
})