<template>
  <div class="page_container">
    <Navigation />
    <router-view />
  </div>
</template>

<script>
import TheNavigation from './components/TheNavigation.vue'

export default {
  name: 'ExpenseApp',
  components: {
    Navigation: TheNavigation,
  },
  data() {
    return {}
  },
}
</script>

<style lang="scss">
.page_container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style>