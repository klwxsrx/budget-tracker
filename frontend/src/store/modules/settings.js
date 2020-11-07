const state = {
  isInitialized: false,
  currency: {
    base: null,
    rates: {},
  },
  firstDayOfPeriod: null,
}

const actions = { // TODO: api request
  ['load']({commit}) {
    return new Promise((resolve) => {
      setTimeout(() => {
        commit('update', {
          currency: {
            base: 'RUB',
            rates: {
              USD: 80.78,
              EUR: 90.34,
            },
          },
          firstDayOfPeriod: 2,
        })
        resolve()
      }, 1000)
    })
  },
}

const mutations = {
  update(state, settings) {
    state.currency.base = settings.currency.base
    state.currency.rates = settings.currency.rates
    state.firstDayOfPeriod = settings.firstDayOfPeriod

    state.isInitialized = true
  },
}

export default {
  namespaced: true,
  state,
  actions,
  mutations,
}