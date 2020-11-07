<template>
  <div class="account-container">
    <ul
      v-if="isLoaded"
      class="account-list"
    >
      <account-list-item
        v-for="account in accounts"
        :key="account.id"
        :account="account"
        class="list-item"
        @selection-changed="onItemSelectionChanged"
      />
    </ul>
    <ul
      v-if="!isLoaded"
      class="account-list"
    >
      <li
        v-for="index in skeletonItemsCount"
        :key="index"
        class="list-item skeleton"
      >
        <div
          :style="[{'flex-grow': Math.random() * skeletonItemsCount + 1}]"
          class="skeleton-title"
        >
          <b-skeleton size="is-medium" />
        </div>
        <div class="skeleton-amount">
          <b-skeleton
            position="is-right"
            size="is-medium"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import AccountListItem from './AccountListItem.vue'

export default {
  name: 'AccountList',
  components: {AccountListItem},
  props: {
    isLoaded: {
      type: Boolean,
      default: false,
    },
    accounts: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      skeletonItemsCount: 5,
    }
  },
  methods: {
    onItemSelectionChanged(id, isSelected) {
      this.$emit('selection-changed', id, isSelected)
    },
  },
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

.account-container {
  background-color: $scheme-main;

  .list-item {
    padding: 16px;

    &.skeleton {
      display: flex;

      .skeleton-title, .skeleton-amount {
        flex-shrink: 1;
      }

      .skeleton-title {
        margin-right: 15%;
      }

      .skeleton-amount {
        flex-grow: 1;
      }
    }
  }
}

@media screen and (min-width: $tablet) {
  .account-container .list-item {
    margin: 10px;
  }
}
</style>