<template>
  <div class="columns is-mobile">
    <AccountList
      :accounts="accounts"
      :is-loaded="isLoaded"
      class="accounts column is-3-desktop is-4-tablet is-12-mobile"
      @selection-changed="changeSelectedAccount"
    />
    <div class="transactions column is-6-desktop is-8-tablet is-12-mobile" />
    <div class="edit-transaction column is-3-desktop is-12-mobile" />
  </div>
</template>

<script>
import {mapState} from 'vuex'
import AccountList from '../components/AccountList.vue'

export default {
  name: 'TransactionsPage',
  components: {
    AccountList: AccountList,
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

<style lang="scss">
@import "~bulma/sass/utilities/initial-variables";

.accounts, .edit-transaction {
  box-shadow: 0 0 2px 0 lightgrey;
  z-index: 1;
}

@media screen and (min-width: $tablet) {
  .columns {
    flex: 1;
  }

  .transactions {
    background-color: #f5f5f5;
  }
}
</style>