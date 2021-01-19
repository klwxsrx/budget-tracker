<template>
  <li
    :class="{selected: account.isSelected}"
    class="item"
    @click="toggleSelected"
  >
    <span class="item-title">{{ account.title }}</span>
    <money-amount
      :amount="account.balance"
      :currency="account.currency"
      class="item-amount"
    />
  </li>
</template>

<script>
import MoneyAmount from './MoneyAmount.vue'

export default {
  name: 'AccountListItem',
  components: {MoneyAmount},
  props: {
    account: {
      type: Object,
      default: () => ({
        id: '',
        title: '',
        balance: 0,
        currency: 'USD',
        isSelected: false,
      }),
    },
  },
  methods: {
    toggleSelected() {
      this.$emit('selection-changed', this.account.id, !this.account.isSelected)
    },
  },
}
</script>

<style lang="scss" scoped>
@import "web/scss/variables";

.item {
  display: flex;
  cursor: pointer;

  &:hover,
  &:hover.selected {
    .item-title,
    .item-amount {
      color: $link;
    }
  }

  &.selected {
    background-color: $primary-light;

    .item-title,
    .item-amount {
      color: $text-strong;
    }
  }

  .item-title {
    flex: 1 1 auto;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    font-size: $size-normal;
    font-weight: $weight-semibold;
    color: $text;
  }

  .item-amount {
    flex: 0 0 auto;
    padding-left: 10px;
    font-size: $size-normal;
    color: $text;
  }
}
</style>