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
        <total-balance-item />
      </b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item
        :active="$route.path === transactionsPath"
        :to="transactionsPath"
        tag="router-link"
      >
        {{ $t('navigation.transactions') }}
      </b-navbar-item>
      <b-navbar-item
        :active="$route.path === settingsPath"
        :to="settingsPath"
        tag="router-link"
      >
        {{ $t('navigation.settings') }}
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
import {routes} from '../router'
import TotalBalanceItem from './TotalBalanceItem.vue'

export default {
  name: 'TheNavigation',
  components: {TotalBalanceItem},
  data() {
    return {
      transactionsPath: routes.transactions.path,
      settingsPath: routes.settings.path,
    }
  },
}
</script>

<style scoped>

</style>