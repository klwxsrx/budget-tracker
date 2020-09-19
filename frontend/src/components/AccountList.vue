<template>
  <div class="account_container">
    <ul
      v-if="isLoaded"
      class="account_list"
    >
      <AccountListItem
        v-for="account in accounts"
        :key="account.id"
        :account="account"
        class="list_item"
        @selection-changed="onItemSelectionChanged"
      />
    </ul>
    <ul
      v-if="!isLoaded"
      class="account_list"
    >
      <li
        v-for="index in skeletonItemsCount"
        :key="index"
        class="list_item skeleton"
      >
        <div
          :style="[{'flex-grow': Math.random() * skeletonItemsCount + 1}]"
          class="skeleton_title"
        >
          <b-skeleton size="is-medium" />
        </div>
        <div class="skeleton_amount">
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
@import "~bulma/sass/utilities/initial-variables";

.account_container .list_item {
  padding: 20px 10px 20px 20px;

  &.skeleton {
    display: flex;

    .skeleton_title, .skeleton_amount {
      flex-shrink: 1;
    }

    .skeleton_title {
      margin-right: 15%;
    }

    .skeleton_amount {
      flex-grow: 1;
    }
  }
}

@media screen and (min-width: $tablet) {
  .account_container {
    margin-top: 10px;
  }
}

@media screen and (max-width: $tablet) {
  .account_container {
    padding-bottom: 0;
  }
}
</style>