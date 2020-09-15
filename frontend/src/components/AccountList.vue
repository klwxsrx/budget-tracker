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
        @selection-changed="changeSelectedAccount"
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
  data() {
    return {
      skeletonItemsCount: 5,
      isLoaded: false,
      accounts: [
        {
          id: 'b17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Наличные',
          balance: 4269,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'a17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Зарплатная, RUB',
          balance: 64000000,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'f17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, RUB',
          balance: 430000,
          currency: 'RUB',
          isSelected: false,
        },
        {
          id: 'c17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, EUR',
          balance: 1400,
          currency: 'EUR',
          isSelected: false,
        },
        {
          id: 'd17aeb44-f459-11ea-adc1-0242ac120002',
          title: 'Дебетовая, USD',
          balance: 17000,
          currency: 'USD',
          isSelected: false,
        },
      ],
    }
  },
  methods: {
    changeSelectedAccount(id, isSelected) {
      this.accounts.forEach(account => {
        account.isSelected = account.id === id && isSelected
      })
      // TODO: change in model
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
      margin-right: 20%;
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