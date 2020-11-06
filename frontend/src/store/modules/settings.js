const state = {
  isInitialized: false,
  currency: {
    base: null,
    rates: {},
  },
}

const actions = { // TODO: api request
  ['load']({commit}) {
    return new Promise((resolve) => {
      setTimeout(() => {
        commit('updateCurrency', {
          base: 'RUB',
          rates: {
            USD: 80.78,
            EUR: 90.34,
          },
        })
        resolve()
      }, 1000)
    })
  },
}

const mutations = {
  updateCurrency(state, settings) {
    state.currency.base = settings.base
    state.currency.rates = settings.rates
  },
}

export default {
  namespaced: true,
  state,
  actions,
  mutations,
}