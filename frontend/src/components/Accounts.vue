<template>
  <div class="account_container">
    <ul>
      <AccountItem v-for="account in accounts"
                   :key="account.id"
                   :balance="account.balance"
                   :currency="account.currency"
                   :title="account.title"
      ></AccountItem>
    </ul>
  </div>
</template>

<script>
import AccountItem from "./AccountItem.vue"

export default {
  name: "Accounts",
  components: {AccountItem},
  data() {
    return {
      accounts: [
        {
          id: "b17aeb44-f459-11ea-adc1-0242ac120002",
          title: "Cash",
          balance: 4269,
          currency: "RUB"
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
</style>