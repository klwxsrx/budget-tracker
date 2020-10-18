<template>
  <li
    :class="{selected: account.isSelected}"
    class="item"
    @click="toggleSelected"
  >
    <span class="item_title">{{ account.title }}</span>
    <MoneyAmount
      :amount="account.balance"
      :currency="account.currency"
      class="item_amount"
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
@import "src/scss/variables";

.item {
  display: flex;
  cursor: pointer;

  &:hover,
  &:hover.selected {
    .item_title,
    .item_amount {
      color: $link;
    }
  }

  &.selected {
    background-color: $scheme-main-ter;

    .item_title,
    .item_amount {
      color: $text-strong;
    }
  }

  .item_title {
    flex: 1 1 auto;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    font-size: $size-normal;
    font-weight: $weight-semibold;
    color: $text;
  }

  .item_amount {
    flex: 0 0 auto;
    padding-left: 10px;
    font-size: $size-normal;
    color: $text;
  }
}
</style>