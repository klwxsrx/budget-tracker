<template>
  <b-field
    v-if="filterItems.length > 0"
    :position="controlsPosition"
    group-multiline
    grouped
  >
    <div
      v-for="item in filterItems"
      :key="item.id"
      class="control"
    >
      <b-tag
        :close-type="item.viewType"
        :size="controlsSize"
        :type="item.viewType"
        attached
        closable
        @close="removeItem(item.id, item.type)"
      >
        {{ item.name }}
      </b-tag>
    </div>
  </b-field>
</template>

<script>
const accountType = 0
const categoryType = 1

function getViewType(type) {
  switch (type) {
  case accountType:
    return 'is-primary'
  case categoryType:
    return 'is-turquoise'
  default:
    return 'is-purple'
  }
}

export default {
  name: 'TransactionFilter',
  props: {
    accountIdFilter: {
      type: String,
      default: null,
    },
    categoryIdFilter: {
      type: String,
      default: null,
    },
  },
  computed: {
    filterItems() {
      let items = []

      const accountId = this.accountIdFilter
      if (accountId !== null) {
        const accountItem = this.getItemByAccountId(accountId)
        if (accountItem !== null) {
          items.push(accountItem)
        }
      }
      return items
    },
    controlsSize() {
      return this.$mq === 'tablet' ? 'is-medium' : 'is-small'
    },
    controlsPosition() {
      return this.$mq === 'tablet' ? 'is-centered' : 'is-left'
    },
  },
  methods: {
    removeItem(id, type) {
      switch (type) {
      case accountType:
        this.$emit('account-filter-removed')
        break
      case categoryType:
        this.$emit('category-filter-removed')
        break
      }
    },
    getItemByAccountId(accountId) {
      const account = this.$store.state.account.items.find(item => {
        return item.id === accountId
      })
      if (account === null) {
        return null
      }
      return {
        id: account.id,
        name: account.title,
        type: accountType,
        viewType: getViewType(accountType),
      }
    },
  },
}
</script>

<style scoped>

</style>