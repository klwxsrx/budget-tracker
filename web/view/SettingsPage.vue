<template>
  <div class="columns">
    <div class="column is-8-desktop is-offset-2-desktop">
      <section>
        <b-field :label="$t('settings.currency')">
          <b-select :placeholder="$t('settings.select-currency')" icon="money-bill-alt">
            <option
                v-for="currency in currencies"
                :key="currency"
                :value="currency"
            >
              {{ $t('currency.' + currency) }}
            </option>
          </b-select>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SettingsPage',
  data() {
    return {
      currencies: ['usd', 'rub'],
    }
  },
}
</script>