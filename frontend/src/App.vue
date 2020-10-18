<template>
  <div class="page-container">
    <TheNavigation />
    <router-view />
  </div>
</template>

<script>
import TheNavigation from './components/TheNavigation.vue'

export default {
  name: 'ExpenseApp',
  components: {
    TheNavigation: TheNavigation,
  },
  data() {
    return {}
  },
}
</script>

<style lang="scss" scoped>
@import "src/scss/variables";

.page-container {
  background-color: $background;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style>