const state = {
  filter: {
    accountId: null,
  },
}

const mutations = {
  filterAccountId(state, value) {
    state.filter.accountId = value
  },
}

export default {
  namespaced: true,
  state,
  mutations,
}