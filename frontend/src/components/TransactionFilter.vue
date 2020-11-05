<template>
  <b-field
    v-if="transactions.length > 0"
    :position="controlsPosition"
    group-multiline
    grouped
  >
    <div
      v-for="transaction in transactions"
      :key="transaction.id"
      class="control"
    >
      <b-tag
        :close-type="transaction.viewType"
        :size="controlsSize"
        :type="transaction.viewType"
        attached
        closable
        @close="removeItem(transaction.id, transaction.type)"
      >
        {{ transaction.name }}
      </b-tag>
    </div>
  </b-field>
</template>

<script>
const accountType = 0
const categoryType = 1
const tagType = 2

function getViewType(type) {
  switch (type) {
  case accountType:
    return 'is-primary'
  case categoryType:
    return 'is-turquoise'
  case tagType:
  default:
    return 'is-purple'
  }
}

export default {
  name: 'TransactionFilter',
  computed: {
    transactions() {
      let items = []

      const accountId = this.$store.state.transaction.filter.accountId
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
        this.$store.commit('transaction/filterAccountId', null)
        break
      default:
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