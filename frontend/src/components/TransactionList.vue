<template>
  <div>
    <div class="shadow-block" />
    <div class="transaction-container">
      <div class="transaction-content">
        <div class="transaction-list-controls">
          <b-field
            class="transaction-date-controls"
            grouped
            position="is-centered"
          >
            <div class="control">
              <b-button
                :size="dateControlsSize"
                icon-left="chevron-left"
                icon-pack="fas"
                type="is-dark"
              />
            </div>
            <div class="control">
              <b-button
                :size="dateControlsSize"
                type="is-dark"
              >
                02 сент. - 02 нояб.
              </b-button>
            </div>
            <div class="control">
              <b-button
                :size="dateControlsSize"
                icon-left="chevron-right"
                icon-pack="fas"
                type="is-dark"
              />
            </div>
          </b-field>
          <b-field
            class="transaction-filter-controls"
            group-multiline
            grouped
            :position="filterControlsPosition"
          >
            <div class="control">
              <b-tag
                attached
                closable
                close-type="is-primary"
                type="is-primary"
                :size="filterControlsSize"
              >
                Наличные
              </b-tag>
            </div>
            <div class="control">
              <b-tag
                attached
                closable
                close-type="is-purple"
                type="is-purple"
                :size="filterControlsSize"
              >
                #питание
              </b-tag>
            </div>
            <div class="control">
              <b-tag
                attached
                closable
                close-type="is-purple"
                type="is-purple"
                :size="filterControlsSize"
              >
                #фастфуд
              </b-tag>
            </div>
          </b-field>
        </div>
        <div class="transaction-list-container columns is-gapless">
          <div class="transaction-list column is-8-desktop is-12-tablet is-12-mobile" />
          <div class="transaction-totals column is-4-desktop is-0-tablet is-0-mobile" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TransactionList',
  computed: {
    dateControlsSize() {
      return this.$mq === 'tablet' ? 'is-medium' : 'is-small'
    },
    filterControlsSize() {
      return this.$mq === 'tablet' ? 'is-medium' : 'is-small'
    },
    filterControlsPosition() {
      return this.$mq === 'tablet' ? 'is-centered' : 'is-left'
    },
  },
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

.transaction-container {
  width: 100%;
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
      background-color: $scheme-main;
      border-top-left-radius: $radius;
      border-top-right-radius: $radius;

      .transaction-list {
        border-top-left-radius: $radius;
      }

      .transaction-totals {
        border-top-right-radius: $radius;
      }
    }
  }
}

@media screen and (min-width: $tablet) {
  .shadow-block {
    box-shadow: 0 0 1px 0 $grey-lighter;
    width: 100%;
    height: 1px;
  }

  .transaction-container {
    height: 100%;

    .transaction-content {
      width: calc(100% - 50px);
      height: calc(100% - 30px);
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
        margin-left: 30px;
      }
    }

    .transaction-list-container {
      box-shadow: 0 0 2px 0 $grey-lighter;
      width: 100%;
      flex-grow: 1;
      margin-top: 10px;
    }
  }
}

@media screen and (max-width: $tablet) {
  .transaction-container .transaction-content {
    margin-top: 30px;
    width: 100%;
    align-items: center;

    .transaction-list-controls {
      flex-direction: column;
      align-items: center;
      width: calc(100% - 32px);

      .transaction-filter-controls {
        margin-top: 30px;
        margin-bottom: 30px;
      }
    }
  }
}

@media screen and (min-width: $desktop) {
  .transaction-container .transaction-list {
    box-shadow: 1px 0 2px -1px $grey-lighter;
  }
}

@media screen and (max-width: $desktop) {
  .transaction-container .transaction-content .transaction-list {
    border-top-right-radius: $radius;
  }
}
</style>