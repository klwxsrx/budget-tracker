const state = {
  isInitialized: false,
  items: [],
}

const actions = {
  ['load']({commit}) {
    return new Promise((resolve) => { // TODO: api request
      setTimeout(() => {
        commit('addItems', [
          {
            id: 'a17aeb44-f459-11ea-adc1-0242ac120002',
            title: 'Наличные',
            balance: 4269,
            currency: 'USD',
            isDeleted: false,
          },
          {
            id: 'b17aeb44-f459-11ea-adc1-0242ac120002',
            title: 'Дебетовая, RUB',
            balance: 30000,
            currency: 'RUB',
            isDeleted: false,
          },
          {
            id: 'c17aeb44-f459-11ea-adc1-0242ac120002',
            title: 'Кредитная, RUB',
            balance: 0,
            currency: 'RUB',
            isDeleted: true,
          },
        ])
        resolve()
      }, 1000)
    })
  },
}

const mutations = {
  addItems(state, items) {
    items.forEach(item => {
      const alreadyExists = state.items.find(existedItem => {
        return item.id === existedItem.id
      })
      if (!alreadyExists) {
        state.items.push(item)
      }
    })
    state.isInitialized = true
  },
}

const getters = {
  notDeletedItems(state) {
    return state.items.filter(item => !item.isDeleted)
  },
}

export default {
  namespaced: true,
  state,
  actions,
  mutations,
  getters,
}