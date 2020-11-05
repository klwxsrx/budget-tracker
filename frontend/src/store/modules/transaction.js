const state = {
  filter: {
    accountId: null,
    tagIds: [],
  },
}

const mutations = {
  filterAccountId(state, value) {
    state.filter.accountId = value
  },
  addTagId(state, value) {
    if (!state.filter.tagIds.includes(value)) {
      state.filter.tagIds.push(value)
    }
  },
  removeTagId(state, value) {
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