const state = {
  filter: { // TODO: delete from model
    accountId: null,
    categoryId: null,
  },
}

const mutations = {
  filterAccountId(state, value) {
    state.filter.accountId = value
  },
  filterCategoryId(state, value) {
    state.filter.categoryId = value
  },
}

export default {
  namespaced: true,
  state,
  mutations,
}