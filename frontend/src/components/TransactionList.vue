<template>
  <div>
    <div class="shadow-block" />
    <div class="transaction-container">
      <div class="transaction-content">
        <div class="transaction-list-controls">
          <div class="transaction-date-controls">
            <b-button
              icon-left="chevron-left"
              icon-pack="fas"
              size="is-small is-dark"
            />
            <b-button size="is-small is-dark">
              02 сент. - 02 нояб.
            </b-button>
            <b-button
              icon-left="chevron-right"
              icon-pack="fas"
              size="is-small is-dark"
            />
          </div>
          <b-field
            class="transaction-filter-controls"
            group-multiline
            grouped
          >
            <div class="control">
              <b-tag
                attached
                closable
                close-type="is-primary"
                type="is-primary"
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
              >
                #фастфуд
              </b-tag>
            </div>
          </b-field>
        </div>
        <div class="transaction-list-container columns is-gapless">
          <div class="transaction-list column is-8-desktop" />
          <div class="transaction-totals column is-4-desktop is-0-tablet is-0-mobile" />
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
      width: 100%;

      .transaction-date-controls {
        flex-shrink: 0;
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

      .transaction-list-controls {
        margin: 0 30px 10px 30px;
        align-items: baseline;

        .transaction-filter-controls
        {
          margin-left: 30px;
          margin-bottom: 10px;
        }
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
}

@media screen and (min-width: $desktop) {
  .transaction-container .transaction-list {
    box-shadow: 1px 0 2px -1px $grey-lighter;
  }
}

@media screen and (max-width: $tablet) {
  .transaction-container {
    .transaction-list-controls {
      flex-direction: column;
      align-items: center;
      margin-top: 30px;
    }

    .transaction-filter-controls {
      margin: 30px 0;
    }
  }
}
</style>