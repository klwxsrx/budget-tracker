import Vue from 'vue';
import VueI18n from 'vue-i18n'

const supportedLocales = [
    'en',
    'ru'
];

Vue.use(VueI18n);
let vueI18n = new VueI18n({
    locale: getUserLocale()
});

function getUserLocale() {
    const userLocale = (navigator.userLanguage || navigator.language).slice(0, 2);
    return supportedLocales.includes(userLocale) ? userLocale : supportedLocales[0];
}

function getCurrentLocale() {
    return getVueI18n().locale;
}

export function getVueI18n() {
    return vueI18n;
}

export default {
    hasCurrencySignBeforeAmount() {
        const hasAfter = [
            'ru'
        ];
        return !hasAfter.includes(getCurrentLocale());
    },

    hasCommaAsMoneyDecimalSeparator() {
        const hasComma = [
            'ru'
        ];
        return hasComma.includes(getCurrentLocale());
    }
}