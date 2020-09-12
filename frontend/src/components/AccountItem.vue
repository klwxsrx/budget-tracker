<template>
  <li class="item">
    <span class="item_title">{{ title }}</span>
    <span class="item_amount">{{ moneyBalance }}</span>
  </li>
</template>

<script>
import moneyBalance from "../mixins/moneyBalance";

export default {
  name: "AccountItem",
  props: {
    title: String,
    balance: Number,
    currency: String
  },
  mixins: [
    moneyBalance
  ]
}
</script>

<style lang="scss" scoped>
.item {
  display: flex;
  padding: 5px;
  margin: 5px;
  font-size: 18px;

  .item_title {
    flex: 1 1 auto;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    margin-left: 10px;
    font-weight: 600;
  }

  .item_amount {
    flex: 0 0 auto;
  }
}
</style>