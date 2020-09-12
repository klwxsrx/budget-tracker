<template>
  <div class="columns is-mobile">
    <Accounts class="column is-3-desktop is-4-tablet is-12-mobile"></Accounts>
    <div class="column is-9-desktop is-8-tablet is-12-mobile"></div>
  </div>
</template>

<script>
import Accounts from "../components/Accounts.vue";

export default {
  name: "Transactions",
  components: {
    Accounts
  }
}
</script>

<style lang="scss">
@import "~bulma/sass/utilities/initial-variables";

@media screen and (min-width: $tablet) {
  .columns {
    flex: 1;
  }
}
</style>