<template>
  <div class="page-container">
    <the-navigation/>
    <router-view/>
  </div>
</template>

<script>
import TheNavigation from './component/TheNavigation.vue'

export default {
  name: 'BudgetTracker',
  components: {TheNavigation},
  data() {
    return {}
  },
}
</script>

<style lang="scss" scoped>
@import "scss/variables";

.page-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: $background;
}
</style>