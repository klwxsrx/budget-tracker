<i18n path="navigation">
{
"en": {
"total_balance": "Total Balance",
"transactions": "Transactions",
"settings": "Settings"
},
"ru": {
"total_balance": "Баланс",
"transactions": "Операции",
"settings": "Настройки"
}
}
</i18n>

<template>
  <b-navbar
      centered
      type="is-dark"
  >
    <template slot="brand">
      <b-navbar-item
          :to="transactionsPath"
          tag="router-link"
      >
        <span class="balance">{{ $t('total_balance') }}</span>
        <MoneyAmount
            :amount="balance"
            :currency="currency"
            class="amount"
        />
      </b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item
          :active="$route.path === transactionsPath"
          :to="transactionsPath"
          tag="router-link"
      >
        {{ $t('transactions') }}
      </b-navbar-item>
      <b-navbar-item
          :active="$route.path === settingsPath"
          :to="settingsPath"
          tag="router-link"
      >
        {{ $t('settings') }}
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
import {routes} from '../router'
import MoneyAmount from './MoneyAmount.vue'

export default {
  name: 'TheNavigation',
  components: {MoneyAmount},
  data() {
    return {
      balance: 42000,
      currency: 'RUB',
      transactionsPath: routes.transactions.path,
      settingsPath: routes.settings.path,
    }
  },
}
</script>

<style lang="scss" scoped>
.balance,
.amount {
  font-size: 18px;
}

.balance {
  margin-left: 18px;
  font-weight: 600;
}

.amount {
  margin-left: 5px;
}
</style>