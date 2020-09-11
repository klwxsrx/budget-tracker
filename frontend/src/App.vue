<template>
  <div>
    <Navigation></Navigation>
    <router-view></router-view>
  </div>
</template>

<script>
import Navigation from "./components/Navigation.vue"

export default {
  name: "ExpenseApp",
  components: {
    Navigation: Navigation
  },
  data() {
    return {};
  }
}
</script>