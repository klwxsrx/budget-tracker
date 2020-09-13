<template>
  <span>{{ moneyAmount }}</span>
</template>

<script>
const currencyMap = {
  'RUB': '₽',
  'USD': '$',
  'EUR': '€'
};

function hasCurrencySign(code) {
  return currencyMap[code] !== undefined;
}

function getCurrencySign(code) {
  return currencyMap[code] || code;
}

function hasCurrencySignBeforeAmount(locale) {
  const currencyMap = {
    'ru': false,
  };
  return currencyMap[locale] !== undefined ? currencyMap[locale] : true;
}

function hasCommaAsDecimalSeparator(locale) {
  const currencyMap = {
    'ru': true
  };
  return currencyMap[locale] !== undefined ? currencyMap[locale] : false;
}

export default {
  name: "MoneyAmount",
  props: {
    amount: Number,
    currency: String
  },
  computed: {
    moneyAmount() {
      const locale = this.$i18n.locale;

      const amountFloat = this.$props.amount / 100;
      let amount = (parseInt(String(amountFloat)) === amountFloat)
          ? String(amountFloat)
          : amountFloat.toFixed(2);
      if (hasCommaAsDecimalSeparator(locale)) {
        amount = amount.replace('.', ',');
      }

      return (hasCurrencySign(this.$props.currency) && hasCurrencySignBeforeAmount(locale))
          ? getCurrencySign(this.$props.currency) + amount
          : amount + getCurrencySign(this.$props.currency);
    }
  }
}
</script>