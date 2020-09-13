<template>
  <div class="page_container">
    <Navigation/>
    <router-view/>
  </div>
</template>

<script>
import Navigation from './components/Navigation.vue'

export default {
  name: 'ExpenseApp',
  components: {
    Navigation: Navigation,
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