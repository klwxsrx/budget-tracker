<template>
  <div class="page-container">
    <TheNavigation />
    <router-view />
  </div>
</template>

<script>
import TheNavigation from './components/TheNavigation.vue'

export default {
  name: 'BudgetTracker',
  components: {
    TheNavigation: TheNavigation,
  },
  data() {
    return {}
  },
  created() {
    this.$store.dispatch('settings/load')
    this.$store.dispatch('account/load')
  },
}
</script>

<style lang="scss" scoped>
@import "web/scss/variables";

.page-container {
  background-color: $background;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style>