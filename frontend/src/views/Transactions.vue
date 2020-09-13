<template>
  <div class="columns is-mobile">
    <Accounts class="accounts column is-3-desktop is-4-tablet is-12-mobile"/>
    <div class="transactions column is-6-desktop is-8-tablet is-12-mobile"/>
    <div class="edit-transaction column is-3-desktop is-12-mobile"/>
  </div>
</template>

<script>
import Accounts from '../components/Accounts.vue'

export default {
  name: 'Transactions',
  components: {
    Accounts,
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