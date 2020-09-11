<template>
  <div class="columns">
    <Accounts class="column is-3"></Accounts>
    <div class="column is-9"></div>
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