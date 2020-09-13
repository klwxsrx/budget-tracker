<template>
  <div class="account_container">
    <ul class="account_list">
      <AccountItem
        v-for="account in accounts"
        :key="account.id"
        :account="account"
        @selection-changed="changeSelectedAccount"
      />
    </ul>
  </div>
</template>

<script>
import AccountItem from './AccountListItem.vue'

export default {
  name: 'AccountList',
  components: {AccountItem},
  data() {
    return {
      accounts: [
        {
          id: 'b17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Наличные',
          balance: 4269,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'a17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Зарплатная, RUB',
          balance: 64000000,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'f17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, RUB',
          balance: 430000,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'c17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, EUR',
          balance: 1400,
          currency: 'EUR',
          isSelected: false,
        },
        {
          id: 'd17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, USD',
          balance: 17000,
          currency: 'USD',
          isSelected: false,
        },
      ],
    }
  },
  methods: {
    changeSelectedAccount(id, isSelected) {
      this.accounts.forEach(account => {
        account.isSelected = account.id === id && isSelected
      })
      // TODO: change in model
    },
  },
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