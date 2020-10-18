<template>
  <div>
    <div class="shadow-block" />
  </div>
</template>

<script>
export default {
  name: 'TransactionList',
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

  .shadow-block {
    box-shadow: 0 0 1px 0 $grey-lighter;
    width: 100%;
    height: 1px;
  }
</style>