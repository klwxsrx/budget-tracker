<i18n path="navigation">
{
"en": {
"total_balance": "Total Balance"
},
"ru": {
"total_balance": "Баланс"
}
}
</i18n>

<template>
  <div>
    <span class="balance">{{ $t('total_balance') }}</span>
    <MoneyAmount
      :amount="balance"
      :currency="currency"
      class="amount"
    />
  </div>
</template>

<script>
import MoneyAmount from './MoneyAmount.vue'

export default {
  name: 'TotalBalanceItem',
  components: {MoneyAmount},
  data() {
    return {
      balance: 42000,
      currency: 'RUB',
    }
  },
}
</script>

<style scoped>

</style>