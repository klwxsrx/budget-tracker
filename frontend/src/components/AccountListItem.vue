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
.item {
  display: flex;
  padding: 20px 1px;
  font-size: 18px;
  cursor: pointer;

  &:hover {
    background-color: #eeeeee;
  }

  &.selected {
    background-color: lightgray;
  }

  .item_title {
    flex: 1 1 auto;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    margin-left: 20px;
    font-weight: 600;
  }

  .item_amount {
    flex: 0 0 auto;
    margin-right: 10px;
  }
}
</style>