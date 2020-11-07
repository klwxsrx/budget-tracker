<template>
  <account-list
    :accounts="accounts"
    :is-loaded="isLoaded"
    @selection-changed="changeSelectedAccount"
  />
</template>

<script>
import {mapState} from 'vuex'
import AccountList from '../components/AccountList.vue'

export default {
  name: 'TransactionPageAccountList',
  components: {AccountList},
  computed: {
    ...mapState({
      isLoaded: state => state.account.isInitialized,
    }),
    accounts() {
      const selectedItemId = this.$store.state.transaction.filter.accountId
      return this.$store.getters['account/notDeletedItems'].map(item => {
        delete item.isDeleted
        this.$set(item, 'isSelected', (item.id === selectedItemId))
        return item
      })
    },
  },
  methods: {
    changeSelectedAccount(id, isSelected) {
      this.$store.commit('transaction/filterAccountId', isSelected ? id : null)
    },
  },
}
</script>

<style scoped>

</style>