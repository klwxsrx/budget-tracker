<template>
  <div>
    <div class="shadow-block"/>
    <div class="transaction-container">
      <div class="transaction-content">
        <div class="transaction-list-controls"></div>
        <div class="transaction-list-container columns is-gapless">
          <div class="transaction-list column is-8-desktop"></div>
          <div class="transaction-totals column is-4-desktop is-0-tablet is-0-mobile"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TransactionList',
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

@media screen and (min-width: $tablet) {
  .shadow-block {
    box-shadow: 0 0 1px 0 $grey-lighter;
    width: 100%;
    height: 1px;
  }

  .transaction-container {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: flex-end;
    justify-content: center;

    .transaction-content {
      width: calc(100% - 50px);
      height: calc(100% - 30px);
      display: flex;
      flex-direction: column;

      .transaction-list-controls {
        width: 100%;
        margin-bottom: 10px;
      }

      .transaction-list-container {
        border-top-left-radius: $radius;
        border-top-right-radius: $radius;
        box-shadow: 0 0 2px 0 $grey-lighter;
        width: 100%;
        height: 100%;

        .transaction-list {
          border-top-left-radius: $radius;
          z-index: 1;
          background-color: $scheme-main;
          height: 100%;
        }

        .transaction-totals {
          border-top-right-radius: $radius;
          background-color: $scheme-main;
          height: 100%;
        }
      }
    }
  }

  @media screen and (min-width: $desktop) {
    .transaction-container .transaction-list {
      box-shadow: 1px 0 2px -1px $grey-lighter;
    }
  }
}
</style>