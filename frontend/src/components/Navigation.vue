<template>
  <b-navbar centered type="is-primary">
    <template slot="brand">
      <b-navbar-item :to="{path: '/'}" tag="router-link">
        <span class="balance">Total Balance</span><span class="amount">10₽</span> <!-- TODO: money mixin-->
      </b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item :active="$route.path === '/'" :to="{path: '/'}" tag="router-link">
        Transactions <!-- TODO: i18n -->
      </b-navbar-item>
      <b-navbar-item :active="$route.path === '/settings'" :to="{path: '/settings'}" tag="router-link">
        Settings
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script>
export default {
  name: "Navigation"
}
</script>

<style scoped>
.balance,
.amount {
  font-size: 18px;
}

.balance {
  margin-left: 8px;
  font-weight: 600;
}

.amount {
  margin-left: 5px;
}
</style>