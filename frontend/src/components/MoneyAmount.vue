<template>
  <span>{{ moneyAmount }}</span>
</template>

<script>
import I18nService from '../services/I18nService'

export default {
  name: 'MoneyAmount',
  props: {
    amount: {
      type: Number,
      default: 0,
    },
    currency: {
      type: String,
      default: 'USD',
    },
  },
  computed: {
    moneyAmount() {
      const amountFloat = this.$props.amount / 100
      return amountFloat.toLocaleString(I18nService.getCurrentLocale(), {
        style: 'currency',
        currency: this.$props.currency,
      })
    },
  },
}
</script>