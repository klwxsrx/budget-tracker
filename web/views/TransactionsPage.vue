<template>
  <div class="columns is-gapless is-tablet">
    <transaction-page-account-list
      :selected-account-id="transactionFilter.accountId"
      class="accounts-block column is-3-desktop is-4-tablet is-12-mobile"
      @account-changed="updateAccountIdFilter"
    />
    <div class="transactions-block column is-9-desktop is-8-tablet is-12-mobile">
      <div class="shadow-block" />
      <div class="transaction-container">
        <div class="transaction-content">
          <div class="transaction-list-controls">
            <transaction-date-range-selector
              :end-date="new Date(2020, 10, 1)"
              :start-date="new Date(2020, 9, 2)"
              class="transaction-date-controls"
            />
            <transaction-filter
              :account-id-filter="transactionFilter.accountId"
              :category-id-filter="transactionFilter.categoryId"
              class="transaction-filter-controls"
              @account-filter-removed="updateAccountIdFilter(null)"
              @category-filter-removed="updateCategoryIdFilter(null)"
            />
          </div>
          <div class="transaction-list-container columns is-gapless">
            <div class="transaction-list column is-8-desktop is-12-tablet is-12-mobile" />
            <div
              v-if="this.$mq !== 'tablet'"
              class="transaction-totals column is-4-desktop is-0-tablet is-0-mobile"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TransactionFilter from '../components/TransactionFilter.vue'
import TransactionPageAccountList from '../components/TransactionPageAccountList.vue'
import TransactionDateRangeSelector from '../components/TransactionDateRangeSelector.vue'

export default {
  name: 'TransactionsPage',
  components: {
    TransactionDateRangeSelector,
    TransactionPageAccountList,
    TransactionFilter,
  },
  data() {
    return {
      transactionFilter: {
        accountId: null,
        categoryId: null,
      },
    }
  },
  methods: {
    updateAccountIdFilter(id, isSelected) {
      this.transactionFilter.accountId = isSelected ? id : null
    },
    updateCategoryIdFilter(id, isSelected) {
      this.transactionFilter.categoryId = isSelected ? id : null
    },
  },
}
</script>

<style lang="scss" scoped>
@import "web/scss/variables";

.columns {
  display: flex;
  flex-grow: 1;

  .accounts-block {
    box-shadow: 0 0 2px 0 $grey-lighter;
    z-index: 1;
  }

  .transactions-block .transaction-container {
    display: flex;
    align-items: flex-end;
    justify-content: center;

    .transaction-content {
      display: flex;
      flex-direction: column;

      .transaction-list-controls {
        display: flex;

        .transaction-date-controls {
          margin-bottom: 0;
        }
      }

      .transaction-list-container {
        box-shadow: 0 0 2px 0 $grey-lighter;
        background-color: $scheme-main;
        width: 100%;

        .transaction-list {
          border-top-left-radius: $radius;
        }

        .transaction-totals {
          border-top-right-radius: $radius;
        }
      }
    }
  }
}

@media screen and (min-width: $tablet) {
  .transactions-block .shadow-block {
    box-shadow: 0 0 1px 0 $grey-lighter;
    width: 100%;
    height: 1px;
  }

  .transactions-block .transaction-container {
    width: 100%;
    height: 100%;

    .transaction-content {
      width: calc(100% - 50px);
      height: calc(100% - 20px);
    }

    .transaction-list-controls {
      padding-left: 20px;
      align-items: baseline;
      width: 100%;

      .transaction-date-controls {
        flex-shrink: 0;
      }

      .transaction-filter-controls
      {
        margin-left: 20px;
      }
    }

    .transaction-list-container {
      border-top-left-radius: $radius;
      border-top-right-radius: $radius;
      flex-grow: 1;
      margin-top: 10px;
    }
  }
}

@media screen and (max-width: $tablet) {
  .columns {
    flex-direction: column;

    .transactions-block {
      display: flex;
      flex-grow: 1;

      .transaction-container {
        flex-grow: 1;

        .transaction-content {
          height: calc(100% - 30px);
          margin-top: 30px;
          width: 100%;
          align-items: center;

          .transaction-list-controls {
            flex-direction: column;
            align-items: center;
            width: calc(100% - 32px);

            .transaction-filter-controls {
              margin-top: 30px;
            }
          }

          .transaction-list-container {
            flex-grow: 1;
            margin-top: 30px;
          }
        }
      }
    }
  }
}

@media screen and (min-width: $desktop) {
  .columns .transactions-block .transaction-container .transaction-list {
    box-shadow: 1px 0 2px -1px $grey-lighter;
  }
}

@media screen and (max-width: $desktop) {
  .columns .transactions-block .transaction-container .transaction-content .transaction-list {
    border-top-right-radius: $radius;
  }
}
</style>