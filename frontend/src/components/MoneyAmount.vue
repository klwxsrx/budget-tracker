<template>
  <span>{{ moneyAmount }}</span>
</template>

<script>
import I18nService from '../services/I18nService'

const currencyMap = {
  'RUB': '₽',
  'USD': '$',
  'EUR': '€',
}

function hasCurrencySign(code) {
  return currencyMap[code] !== undefined
}

function getCurrencySign(code) {
  return currencyMap[code] || code
}

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
      let amount = (parseInt(String(amountFloat)) === amountFloat)
          ? String(amountFloat)
          : amountFloat.toFixed(2)
      if (I18nService.hasCommaAsMoneyDecimalSeparator()) {
        amount = amount.replace('.', ',')
      }

      return (hasCurrencySign(this.$props.currency) && I18nService.hasCurrencySignBeforeAmount())
          ? getCurrencySign(this.$props.currency) + amount
          : amount + getCurrencySign(this.$props.currency)
    },
  },
}
</script>