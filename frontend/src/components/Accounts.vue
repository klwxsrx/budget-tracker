<template>
  <div class="account_container">
    <ul>
      <li v-for="account in accounts" v-bind:key="account.id" class="item">
        <span class="account_title">{{ account.title }}</span>
        <span class="account_amount">{{ account.currentBalance }}₽</span> <!-- TODO: money mixin -->
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: "Accounts",
  data() {
    return {
      accounts: [
        {
          id: "b17aeb44-f459-11ea-adc1-0242ac120002",
          title: "Cash",
          currentBalance: 42.69,
          currency: "USD"
        }
      ]
    };
  }
}
</script>

<style lang="scss" scoped>
.account_container {
  background-color: #f5f5f5;
}

.item {
  display: flex;
  padding: 5px;
  margin: 5px;

  .account_title {
    flex: 1 1 auto;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    margin-left: 10px;
    font-size: 18px;
    font-weight: 600;
  }

  .account_amount {
    flex: 0 0 auto;
    margin-top: 1px;
  }
}
</style>