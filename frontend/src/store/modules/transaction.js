const state = {
  filter: {
    accountId: null,
    categoryId: null,
    tagIds: [],
  },
}

const mutations = {
  filterAccountId(state, value) {
    state.filter.accountId = value
  },
  filterCategoryId(state, value) {
    state.filter.categoryId = value
  },
  addTagIdToFilter(state, value) {
    if (!state.filter.tagIds.includes(value)) {
      state.filter.tagIds.push(value)
    }
  },
  removeTagIdFromFilter(state, value) {
    const index = state.filter.tagIds.indexOf(value)
    if (index > -1) {
      state.filter.tagIds.splice(index, 1)
    }
  },
}

export default {
  namespaced: true,
  state,
  mutations,
}