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

export default {
    props: {
        balance: Number,
        currency: String
    },
    computed: {
        moneyBalance() {
            const currencySignBeforeAmount = false; // TODO: i18n
            const commaIsDecimalSeparator = true;
            const balanceFloat = this.$props.balance / 100;

            let balance = (parseInt(String(balanceFloat)) === balanceFloat)
                ? balanceFloat
                : balanceFloat.toFixed(2);
            if (commaIsDecimalSeparator) {
                balance = balance.replace('.', ',');
            }

            return (hasCurrencySign(this.$props.currency) && currencySignBeforeAmount)
                ? getCurrencySign(this.$props.currency) + balance
                : balance + getCurrencySign(this.$props.currency);
        }
    }
};