<template>
  <div class="total-balance">
    <span class="balance">{{ $t('balance.total_balance') }} </span>
    <money-amount
      v-if="isLoaded"
      :amount="totalBalance.balance"
      :currency="totalBalance.currency"
      :prefix="hasOtherCurrencies ? '~' : ''"
      class="amount"
    />
    <b-icon
      v-else
      custom-class="fa-spin"
      icon="spinner"
      pack="fas"
    />
  </div>
</template>

<script>
import MoneyAmount from './MoneyAmount.vue'
import {mapState} from 'vuex'

export default {
  name: 'TotalBalanceItem',
  components: {MoneyAmount},
  computed: {
    ...mapState({
      isLoaded: state => state.account.isInitialized && state.settings.isInitialized,
    }),
    hasOtherCurrencies() {
      const baseCurrency = this.$store.state.settings.currency.base
      return Boolean(this.$store.getters['account/notDeletedItems'].find(item => item.currency !== baseCurrency))
    },
    totalBalance() {
      const currencySettings = this.$store.state.settings.currency
      return {
        balance: this.calculateBalance(
          this.$store.getters['account/notDeletedItems'],
          currencySettings.base,
          currencySettings.rates,
        ),
        currency: currencySettings.base,
      }
    },
  },
  methods: {
    calculateBalance(accounts, base, rates) {
      let result = 0
      accounts.forEach(account => {
        const currency = account.currency
        const hasBaseCurrency = currency === base

        if (!hasBaseCurrency && !rates[currency]) {
          console.error('account currency ' + currency +' is unknown!')
          return null
        }

        result += hasBaseCurrency ? account.balance : (account.balance * rates[currency])
      })

      return result
    },
  },
}
</script>

<style lang="scss" scoped>
@import "web/scss/variables";

.total-balance {
  min-width: 200px;

  .balance,
  .amount {
    font-size: $size-normal;
  }

  .balance {
    margin-left: 14px;
    font-weight: $weight-semibold;
  }

  .amount {
    margin-left: 5px;
  }
}

@media screen and (max-width: $tablet) {
  .total-balance .balance {
    margin-left: 4px;
  }
}
</style>