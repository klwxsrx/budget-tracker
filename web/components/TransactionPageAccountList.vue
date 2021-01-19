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
  props: {
    selectedAccountId: {
      type: String,
      default: null,
    },
  },
  computed: {
    ...mapState({
      isLoaded: state => state.account.isInitialized,
    }),
    accounts() {
      const selectedItemId = this.$props.selectedAccountId
      return this.$store.getters['account/notDeletedItems'].map(item => {
        delete item.isDeleted
        this.$set(item, 'isSelected', (item.id === selectedItemId))
        return item
      })
    },
  },
  methods: {
    changeSelectedAccount(id, isSelected) {
      this.$emit('account-changed', id, isSelected)
    },
  },
}
</script>

<style scoped>

</style>