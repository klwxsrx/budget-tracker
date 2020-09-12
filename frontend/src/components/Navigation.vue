<template>
  <b-navbar centered type="is-primary">
    <template slot="brand">
      <b-navbar-item :to="transactionsPath" tag="router-link">
        <span class="balance">Total Balance</span><span class="amount">{{ moneyBalance }}</span>
      </b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item :active="$route.path === transactionsPath" :to="transactionsPath" tag="router-link">
        Transactions <!-- TODO: i18n -->
      </b-navbar-item>
      <b-navbar-item :active="$route.path === settingsPath" :to="settingsPath" tag="router-link">
        Settings
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
import {routes} from "../router";
import moneyBalance from "../mixins/moneyBalance";

export default {
  name: "Navigation",
  props: {
    balance: Number,
    currency: String
  },
  mixins: [moneyBalance],
  data() {
    return {
      transactionsPath: routes.transactions.path,
      settingsPath: routes.settings.path
    };
  }
}
</script>

<style lang="scss" scoped>
@import "~bulma/sass/utilities/initial-variables";

.balance,
.amount {
  font-size: $size-5;
}

.balance {
  margin-left: 8px;
  font-weight: $weight-semibold;
}

.amount {
  margin-left: 5px;
}
</style>