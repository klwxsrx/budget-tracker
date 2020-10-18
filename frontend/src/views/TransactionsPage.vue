<template>
  <div class="columns is-gapless is-tablet">
    <AccountList
      :accounts="accounts"
      :is-loaded="isLoaded"
      class="accounts column is-3-desktop is-4-tablet is-12-mobile"
      @selection-changed="changeSelectedAccount"
    />
    <TransactionList class="transactions column is-9-desktop is-8-tablet is-12-mobile" />
  </div>
</template>

<script>
import {mapState} from 'vuex'
import AccountList from '../components/AccountList.vue'
import TransactionList from '../components/TransactionList.vue'

export default {
  name: 'TransactionsPage',
  components: {
    TransactionList,
    AccountList,
  },
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
  created() {
    this.$store.dispatch('account/loadAccounts')
  },
  methods: {
    changeSelectedAccount(id, isSelected) {
      if (isSelected) {
        this.$store.commit('transaction/filterAccountId', id)
      } else {
        this.$store.commit('transaction/filterAccountId', null)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

.columns {
  .accounts {
    box-shadow: 0 0 2px 0 $grey-lighter;
    z-index: 1;
  }
}

@media screen and (min-width: $tablet) {
  .columns {
    flex: 1;
  }
}
</style>