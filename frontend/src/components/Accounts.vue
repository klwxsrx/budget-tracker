<template>
  <div class="account_container">
    <ul class="account_list">
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
        },
        {
          id: "a17aeb44-f459-11ea-adc1-0242ac120002",
          title: "Debt Card",
          balance: 6745,
          currency: "RUB"
        }
      ]
    };
  }
}
</script>

<style lang="scss" scoped>
@import "~bulma/sass/utilities/initial-variables";

@media screen and (min-width: $tablet) {
  .account_container {
    margin-top: 10px;
  }
}

@media screen and (max-width: $tablet) {
  .account_container {
    padding-bottom: 0;
  }
}
</style>