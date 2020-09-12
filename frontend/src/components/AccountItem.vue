<template>
  <li :class="{selected: isSelected}" class="item" @click="toggleSelection">
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
    currency: String,
  },
  mixins: [
    moneyBalance
  ],
  methods: {
    toggleSelection() {
      this.isSelected = !this.isSelected;
    }
  },
  data() {
    return {
      isSelected: false
    }
  }
}
</script>

<style lang="scss" scoped>
.item {
  display: flex;
  padding: 10px;
  font-size: 18px;
  cursor: pointer;

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
  }
}
</style>